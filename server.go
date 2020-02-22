package main

import (
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type Scope struct {
	name string
	identifier string
}

func main() {
	LoadConfig()

	InitStore()

	r := mux.NewRouter()

	commonMiddleware := negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.HandlerFunc(PathSanitizer))

	r.NotFoundHandler = http.HandlerFunc(NotFoundMiddleware)
	r.PathPrefix(javascriptDirectory).Handler(commonMiddleware.With(
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.PathPrefix(mediaDirectory).Handler(commonMiddleware.With(
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.Path("/callback").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(CallbackHandler))))

	r.Path("/login").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(LoginHandler))))

	r.PathPrefix(protectedDirectory).Handler(commonMiddleware.With(
		negroni.HandlerFunc(ProtectedResourceMiddleware),
		negroni.Wrap(http.FileServer(http.Dir(staticDirectory)))))

	r.Path("/cookie").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(SetCookie))))

	r.PathPrefix("/").Handler(commonMiddleware.With(
		negroni.HandlerFunc(RedirectHomeMiddleware),
		negroni.Wrap(http.HandlerFunc(IndexHandler(staticDirectory + indexPage)))))

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serveHost),
	}

	dir := cacheDir()
	if dir != "" {
		certManager.Cache = autocert.DirCache(dir)
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
		fmt.Println("will serve " + entrypoint)
		http.ServeFile(w, r, entrypoint)
	}
	return fn
}

func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Printf("%s: %s%s - %s\n", time.Now().Format(time.RFC3339), r.Host, r.URL, r.Method)
	next(rw, r)
}

func PathSanitizer(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if strings.Contains(r.URL.String(), "..") {
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

// cacheDir makes a consistent cache directory inside /tmp. Returns "" on error.
func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "cache-golang-autocert-"+u.Username)
		if err := os.MkdirAll(dir, 0700); err == nil {
			return dir
		}
	}
	return ""
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
	cookie, _ := r.Cookie("protected_cookie")
	if cookie != nil {
		fmt.Println(cookie.Raw)
		next(rw, r)
	}
	rw.WriteHeader(403)
	rw.Write([]byte("<head>Forbidden</head>"))
}

func OpenResource(path string, resource string) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("serving resource")
		http.ServeFile(rw, r, path + resource)
	}
}