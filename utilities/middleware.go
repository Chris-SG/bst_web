package utilities

import (
	"github.com/urfave/negroni"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	commonMiddleware *negroni.Negroni
	protectionMiddleware *negroni.Negroni
	cachingMiddleware *negroni.Negroni

	logger *negroni.Logger
)

func PrepareMiddleware() {
	logger = negroni.NewLogger()
	// MIDDLEWARE DEFINITIONS
	commonMiddleware = negroni.New(
		negroni.HandlerFunc(logger.ServeHTTP),
		negroni.HandlerFunc(PathSanitizer),
		negroni.HandlerFunc(RefreshJwt),
		negroni.HandlerFunc(LogoutIfExpired))

	protectionMiddleware = negroni.New(
		negroni.HandlerFunc(ProtectedResourceMiddleware))

	cachingMiddleware = negroni.New(
		negroni.HandlerFunc(FileCacher))
}

func GetCommonMiddleware() *negroni.Negroni {
	return commonMiddleware
}

func GetProtectionMiddleware() *negroni.Negroni {
	return protectionMiddleware
}

func GetCachingMiddleware() *negroni.Negroni {
	return cachingMiddleware
}


func FileCacher(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Cache-Control", "max-age=3600")
	upath := r.URL.Path
	path.Clean(upath)
	root := http.Dir(StaticDirectory)
	fs, _ := root.Open(upath)

	var modTime time.Time
	fi, err := fs.Stat()
	if err != nil {
		modTime = fi.ModTime()
	} else {
		modTime = time.Now()
	}
	etag := "\"" + upath + modTime.String() + "\""
	rw.Header().Set("Etag", etag)

	next(rw, r)
}

func ProtectedResourceMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["profile"] == nil {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"status":"bad","message":"you are not currently logged in."}`))
		return
	}

	profile := session.Values["profile"].(map[string]interface{})
	expTime := time.Unix(int64(profile["exp"].(float64)), 0)
	if expTime.Unix() < time.Now().Unix() {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"status":"bad","message":"your session has expired."}`))
		return
	}

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
