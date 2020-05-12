package main

import (
	"bst_web/utilities"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	utilities.LoadConfig()
	utilities.PrepareMiddleware()

	utilities.InitStore()
	utilities.InitClient()

	r := mux.NewRouter()


	r.NotFoundHandler = http.HandlerFunc(utilities.NotFoundMiddleware)

	r.Path("/{path:.*\\.js$}").Handler(utilities.GetCommonMiddleware().With(
		negroni.HandlerFunc(SetContentType("application/javascript")),
		negroni.Wrap(utilities.GetCachingMiddleware().With(
			negroni.Wrap(http.FileServer(http.Dir(utilities.StaticDirectory)))))))

	// SUB-ROUTERS
	r.PathPrefix("/external").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(CreateExternalRouters("", nil))))

	r.PathPrefix("/user").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(utilities.GetProtectionMiddleware().With(
			negroni.Wrap(UserRouter())))))

	r.PathPrefix("/ddr").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(utilities.GetProtectionMiddleware().With(
		negroni.Wrap(DdrRouter())))))

	r.PathPrefix("/drs").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(utilities.GetProtectionMiddleware().With(
		negroni.Wrap(DrsRouter())))))

	AttachAuthRoutes(r)

	r.Path("/whoami").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(http.HandlerFunc(WhoAmI)))).Methods(http.MethodGet)

	r.Path("/token").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			session, _ := utilities.Store.Get(r, "auth-session")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprint(session)))
	})))).Methods(http.MethodGet)

	// FILESERVERS

	r.PathPrefix(utilities.MediaDirectory).Handler(utilities.GetCommonMiddleware().With(
		negroni.HandlerFunc(SetMediaContentType),
		negroni.Wrap(http.FileServer(http.Dir(utilities.StaticDirectory)))))

	r.PathPrefix(utilities.CssDirectory).Handler(utilities.GetCommonMiddleware().With(
		negroni.HandlerFunc(SetContentType("text/css")),
		negroni.Wrap(http.FileServer(http.Dir(utilities.StaticDirectory)))))

	r.PathPrefix("/").Handler(utilities.GetCommonMiddleware().With(
		negroni.HandlerFunc(utilities.RedirectHomeMiddleware),
		negroni.Wrap(http.HandlerFunc(IndexHandler(utilities.StaticDirectory+utilities.IndexPage)))))

	var certManager *autocert.Manager

	certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(utilities.ServeHost),
		Cache: autocert.DirCache("./cert_cache"),
	}

	srv := &http.Server{
		Handler:           r,
		Addr:		":" + utilities.ServePort,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 90 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		h := certManager.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	log.Fatal(srv.ListenAndServeTLS("", ""))
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return fn
}

func SetContentType(contentType string) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.Header().Set("Content-Type", contentType)
		next(rw, r)
	}
}

func SetMediaContentType(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	next(rw, r)
}

func OpenResource(path string, resource string) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("serving resource")
		http.ServeFile(rw, r, path + resource)
	}
}

func WhoAmI(rw http.ResponseWriter, r *http.Request) {
	type CacheableData struct {
		Id int `json:"id"`
		Nickname string `json:"nickname"`
		Public bool `json:"public"`
	}

	session, err := utilities.Store.Get(r, "auth-session")
	if err != nil || session == nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(""))
		return
	}
	if _, ok := session.Values["access_token"]; !ok {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(""))
		return
	}

	profileMap, ok := session.Values["profile"].(map[string]interface{})
	if ok {
		sub, ok := profileMap["sub"].(string)
		if ok {
			sub = strings.ToLower(sub)
			cacheResult := utilities.GetCacheValue("users", sub)
			if cacheResult == nil {
				LoadUserCache(sub)
				cacheResult = utilities.GetCacheValue("users", sub)
				if cacheResult == nil {
					rw.WriteHeader(http.StatusOK)
					rw.Write([]byte(""))
					return
				}
			}

			userCache := cacheResult.(CacheableData)
			user, _ := json.Marshal(userCache)
			rw.WriteHeader(http.StatusOK)
			rw.Write(user)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(""))
	return
}

func LoadUserCache(user string) bool {
	type CacheableData struct {
		Id int `json:"id"`
		Nickname string `json:"nickname"`
		Public bool `json:"public"`
	}

	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "cache")
	uri.Query().Set("user", user)

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
	}

	res, err := utilities.GetClient().Do(req)
	if err != nil {
		return false
	}

	cacheData := CacheableData{}
	json.NewDecoder(res.Body).Decode(&cacheData)

	return utilities.SetCacheValue("users", user, cacheData)
}