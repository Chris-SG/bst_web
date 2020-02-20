package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)
var (
	Store *sessions.FilesystemStore
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func InitStore() error {
	Store = sessions.NewFilesystemStore("", []byte(fileStoreKey))
	gob.Register(map[string]interface{}{})
	return nil
}

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
		RedirectURL:  "https://" + serveHost + callbackResourcePath,
		Endpoint: 	  provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

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
	session.Values["profile"] = profile
	err = session.Save(r, rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(rw, r, "/user", http.StatusSeeOther)
}

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
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
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

	http.Redirect(rw, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}