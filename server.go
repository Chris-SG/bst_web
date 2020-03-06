package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var (
	commonMiddleware *negroni.Negroni
	protectionMiddleware *negroni.Negroni
)

func main() {
	LoadConfig()

	InitStore()
	InitClient()

	r := mux.NewRouter()
	logger := negroni.NewLogger()

	// MIDDLEWARE DEFINITIONS
	commonMiddleware = negroni.New(
		negroni.HandlerFunc(logger.ServeHTTP),
		negroni.HandlerFunc(PathSanitizer),
		negroni.HandlerFunc(RefreshJwt),
		negroni.HandlerFunc(LogoutIfExpired))

	protectionMiddleware = negroni.New(
		negroni.HandlerFunc(ProtectedResourceMiddleware))


	r.NotFoundHandler = http.HandlerFunc(NotFoundMiddleware)

	// SUB-ROUTERS
	r.PathPrefix("/external").Handler(commonMiddleware.With(
		negroni.Wrap(CreateExternalRouters(nil))))

	r.PathPrefix("/user").Handler(commonMiddleware.With(
		negroni.Wrap(protectionMiddleware.With(
			negroni.Wrap(UserRouter())))))

	AttachAuthRoutes(r)

	// FILESERVERS
	r.PathPrefix(javascriptDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(SetContentType("application/javascript")),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.PathPrefix(mediaDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(SetMediaContentType),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.PathPrefix(cssDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(SetContentType("text/css")),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.PathPrefix("/").Handler(commonMiddleware.With(
		negroni.HandlerFunc(RedirectHomeMiddleware),
		negroni.Wrap(http.HandlerFunc(IndexHandler(staticDirectory + indexPage)))))

	var certManager *autocert.Manager

	certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serveHost),
		Cache: autocert.DirCache("./cert_cache"),
	}

	srv := &http.Server{
		Handler:           r,
		Addr:		":" + servePort,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
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
		header := LoadHeader(r)
		footer := LoadFooter()
		t, _:= template.ParseFiles(entrypoint)
		replace := struct {
			Header string
			Footer string
		} {
			header,
			footer,
		}

		t.Execute(w, replace)
		//http.ServeFile(w, r, entrypoint)
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

func PathSanitizer(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if strings.Contains(r.URL.String(), "..") ||
	   strings.Contains(r.URL.String(), "./") {
		NotFoundMiddleware(rw, r)
		return
	}

	next(rw, r)
}

func RedirectHomeMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path != "/" {
		NotFoundMiddleware(rw, r)
		return
	}
	next(rw, r)
}

func NotFoundMiddleware(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "https://" + r.Host, 301)
}

func ProtectedResourceMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["profile"] == nil {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("you are not currently logged in."))
		return
	}

	profile := session.Values["profile"].(map[string]interface{})
	expTime := time.Unix(int64(profile["exp"].(float64)), 0)
	if expTime.Unix() < time.Now().Unix() {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("your session has expired."))
		return
	}

	next(rw, r)
}

func OpenResource(path string, resource string) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("serving resource")
		http.ServeFile(rw, r, path + resource)
	}
}