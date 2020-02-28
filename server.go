package main

import (
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
	"text/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	LoadConfig()

	InitStore()
	InitClient()

	r := mux.NewRouter()

	logger := negroni.NewLogger()

	commonMiddleware := negroni.New(
		negroni.HandlerFunc(logger.ServeHTTP),
		negroni.HandlerFunc(PathSanitizer))

	r.NotFoundHandler = http.HandlerFunc(NotFoundMiddleware)
	r.PathPrefix(javascriptDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(SetContentType("application/javascript")),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.PathPrefix(mediaDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(SetMediaContentType),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.PathPrefix(cssDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(SetContentType("text/css")),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.Path("/callback").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(CallbackHandler))))

	r.Path("/login").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(LoginHandler))))

	r.Path("/logout").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(LogoutHandler))))

	r.PathPrefix("/user").Handler(commonMiddleware.With(
		negroni.HandlerFunc(ProtectedResourceMiddleware),
		negroni.Wrap(http.FileServer(http.Dir(protectedDirectory)))))

	ddrRouter := mux.NewRouter().PathPrefix("/ddr").Subrouter()
	ddrRouter.HandleFunc("/songs", DDRSongs)
	ddrRouter.HandleFunc("/songs/{id}", DDRSongsId)

	r.PathPrefix("/ddr").Handler(commonMiddleware.With(
		negroni.Wrap(ddrRouter)))

	r.PathPrefix("/ajax").Handler(commonMiddleware.With(
		negroni.Wrap(AjaxRouter())))

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

func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Printf("%s: %s%s - %s\n", time.Now().Format(time.RFC3339), r.Host, r.URL, r.Method)
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

func SetCookie(rw http.ResponseWriter, r *http.Request) {
	expireTime := time.Now().AddDate(0,0, 1)
	uuid := uuid.New().String()
	cookie := &http.Cookie{
		Name:       "protected_cookie",
		Value:      uuid,
		Domain:		serveHost,
		Path:       "/",
		Expires:    expireTime,
		RawExpires: expireTime.Format(time.UnixDate),
		MaxAge:     86400,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteDefaultMode,
		Raw:        "protected_cookie=" + uuid,
		Unparsed:   []string{"protected_cookie" + uuid},
	}
	fmt.Println(cookie)
	//validCookies := append(validCookies, uuid)
	http.SetCookie(rw, cookie)
}

func ProtectedResourceMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
	rw.Write([]byte(fmt.Sprintf("<head></head><body>%s</body>", session)))
}

func OpenResource(path string, resource string) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("serving resource")
		http.ServeFile(rw, r, path + resource)
	}
}