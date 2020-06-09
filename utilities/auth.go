package utilities

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	bst_models "github.com/chris-sg/bst_server_models"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	Store *sessions.FilesystemStore
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

// InitStore will ensure a store for auth data exists.
func InitStore() error {
	Store = sessions.NewFilesystemStore("./store", []byte(fileStoreKey))
	gob.Register(map[string]interface{}{})
	return nil
}

// NewAuthenticator will provide an authenticator to be used against
// the designated authorization server.
func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, authClientIssuer)
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     authClientId,
		ClientSecret: authClientSecret,
		RedirectURL:  "https://" + ServeHost + callbackResourcePath,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "offline_access", "database"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

// CallbackHandler handles the token exchange part of the authorzization flow.
func CallbackHandler(rw http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(rw, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(rw, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	refreshToken, ok := token.Extra("refresh_token").(string)
	if !ok {
		http.Error(rw, "No refresh_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: authClientId,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Error(rw, "Failed to verify ID Token: " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = refreshToken
	session.Values["profile"] = profile
	err = session.Save(r, rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(rw, r, "/user", http.StatusSeeOther)
}

// LoginHandler will create a session for the user and initiate the
// login flow.
func LoginHandler(rw http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := Store.Get(r, "auth-session")
	session.Values["state"] = state
	err = session.Save(r, rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, authenticator.Config.AuthCodeURL(state, oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
}

// Logout will remove a users login state.
func LogoutHandler(rw http.ResponseWriter, r *http.Request) {
	logoutUrl, err := url.Parse(authClientIssuer)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" +  r.Host)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", authClientId)
	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(rw, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}

// RefreshJwt will send a refresh request to the designated authorization server
// in case expiry is near.
func RefreshJwt(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		next(rw, r)
		return
	}

	if session.Values["profile"] == nil {
		next(rw, r)
		return
	}

	profile := session.Values["profile"].(map[string]interface{})
	expTime := time.Unix(int64(profile["exp"].(float64)), 0)
	iat := time.Unix(int64(profile["iat"].(float64)), 0)
	buffer := (expTime.Unix() - iat.Unix()) / 10
	if expTime.Unix() < time.Now().Add(time.Second * time.Duration(buffer)).Unix() {
		refreshToken := session.Values["refresh_token"].(string)
		authEndpoint, err := url.Parse(authClientIssuer + "oauth/token")

		data := url.Values{
			"grant_type": {"refresh_token"},
			"client_id": {authClientId},
			"client_secret": {authClientSecret},
			"refresh_token": {refreshToken},
		}

		req, _ := http.NewRequest(http.MethodPost, authEndpoint.String(), strings.NewReader(data.Encode()))
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
		req.Header.Add("content-length", strconv.Itoa(len(data.Encode())))
		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(body, &responseMap)
		if err != nil {
			fmt.Println(err)
			next(rw, r)
			return
		}

		rawIDToken, ok := responseMap["id_token"].(string)
		if !ok {
			http.Error(rw, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: authClientId,
		}

		authenticator, err := NewAuthenticator()
		if err != nil {
			fmt.Println(err)
			next(rw, r)
			return
		}

		idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

		if err != nil {
			fmt.Println(err)
			next(rw, r)
			return
		}

		var updatedProfile map[string]interface{}
		if err := idToken.Claims(&updatedProfile); err != nil {
			fmt.Println(err)
			next(rw, r)
			return
		}

		session.Values["id_token"] = rawIDToken
		session.Values["access_token"] = responseMap["access_token"]
		session.Values["refresh_token"] = refreshToken
		session.Values["profile"] = updatedProfile
		err = session.Save(r, rw)
		if err != nil {
			fmt.Println(err)
			next(rw, r)
			return
		}
	}

	next(rw, r)
}

// LogoutIfExpired will ensure the user is logged out if the token happens to
// expire without being successfully renewed.
func LogoutIfExpired(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		next(rw, r)
		return
	}

	if session.Values["profile"] == nil {
		next(rw, r)
		return
	}

	profile := session.Values["profile"].(map[string]interface{})
	expTime := time.Unix(int64(profile["exp"].(float64)), 0)
	if expTime.Unix() < time.Now().Unix() {
		cookie := &http.Cookie {
			Name:    "auth-session",
			Value:   "",
			Expires: time.Unix(0, 0),
			Domain:  ServeHost,
			Path:    "/",
		}
		http.SetCookie(rw, cookie)
	}
	next(rw, r)
}

// TokenRequest retrieves the users id token.
func TokenForRequest(r *http.Request) (token string, err bst_models.Error) {
	err = bst_models.ErrorOK
	session, e := Store.Get(r, "auth-session")
	if e != nil {
		err = bst_models.ErrorJwt
		return
	}

	if session.Values["id_token"] == nil {
		err = bst_models.ErrorJwt
		return
	}

	token = session.Values["id_token"].(string)
	return
}

func ProfileForRequest(r *http.Request) (profile map[string]interface{}, err bst_models.Error) {
	err = bst_models.ErrorOK
	session, e := Store.Get(r, "auth-session")
	if e != nil {
		err = bst_models.ErrorJwt
		return
	}

	if _, ok := session.Values["access_token"]; !ok {
		err = bst_models.ErrorJwt
		return
	}

	profile, ok := session.Values["profile"].(map[string]interface{})
	if !ok {
		err = bst_models.ErrorJwtProfile
	}
	return
}