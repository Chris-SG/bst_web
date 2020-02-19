package main

import (
	"crypto/tls"
	"flag"
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
	"time"
)

type Scope struct {
	name string
	identifier string
}

/*var (
	validCookies []string
	protectedScopeRequired map[string][]Scope
)*/

var (
	clientId string
	clientSecret string
	issuer string
	host string
)

// sudo ./bst_web -static="./dist" -entry="./dist/index.html" -clientid="" -clientsecret="" -issuer="" -host="abc.com"

func main() {
	var entry string
	var port string
	var jsdir string
	var protectedDirectory string
	var static string

	flag.StringVar(&entry, "entry", "./index.html", "the entrypoint to serve.")
	flag.StringVar(&static, "static", "./dist", "the directory to serve static files from.")
	flag.StringVar(&jsdir, "jsdir", "./dist/js", "the directory to serve js files from")
	flag.StringVar(&port, "port", "8000", "the `port` to listen on.")
	flag.StringVar(&protectedDirectory, "protecteddirectory", "./dist/protected", "the directory containing protected files")
	flag.StringVar(&clientId, "clientid", "", "client id for oauth2 client")
	flag.StringVar(&clientSecret, "clientsecret", "", "client secret for oauth2 client")
	flag.StringVar(&issuer, "issuer", "", "issuer for oauth2 client")
	flag.StringVar(&host, "host", "", "host")
	flag.Parse()

	InitStore()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(NotFoundMiddleware)
	r.PathPrefix("/js").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.Wrap(http.FileServer(http.Dir(static)))))

	r.PathPrefix("/img").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.Wrap(http.FileServer(http.Dir(static)))))

	r.Path("/callback").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.Wrap(http.HandlerFunc(CallbackHandler))))

	r.Path("/login").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.Wrap(http.HandlerFunc(LoginHandler))))

	r.Path("/protected").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.HandlerFunc(ProtectedResourceMiddleware),
		negroni.Wrap(http.HandlerFunc(OpenResource(static, "protected.html")))))

	r.Path("/cookie").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.Wrap(http.HandlerFunc(SetCookie))))

	r.PathPrefix("/").Handler(negroni.New(
		negroni.HandlerFunc(LoggingMiddleware),
		negroni.HandlerFunc(RedirectHomeMiddleware),
		negroni.Wrap(http.HandlerFunc(IndexHandler(entry)))))

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
	}

	dir := cacheDir()
	if dir != "" {
		certManager.Cache = autocert.DirCache(dir)
	}

	srv := &http.Server{
		Handler:           r,
		Addr:		":443",
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
		http.ServeFile(w, r, entrypoint)
	}
	return fn
}

func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Printf("%s: %s%s - %s\n", time.Now().Format(time.RFC3339), r.Host, r.URL, r.Method)
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
		Domain:		host,
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