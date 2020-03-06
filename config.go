package main

import "flag"

var (
	staticDirectory string
	indexPage string
	notFoundPage string
	cssDirectory string
	javascriptDirectory string
	mediaDirectory string
	publicDirectory string
	protectedDirectory string

	authClientId string
	authClientSecret string
	authClientIssuer string
	authClientAudience string
	callbackResourcePath string

	fileStoreKey string

	serveHost string
	servePort string

	bstApi string
	bstApiBase string
)

// LoadConfig populates general configuration values to be used with the program.
func LoadConfig() {
	flag.StringVar(&staticDirectory, "static", "./dist", "the directory containing all static files.")
	flag.StringVar(&indexPage, "index", "/index.html", "the location of the index page, relative to the `static` directory.")
	flag.StringVar(&notFoundPage, "404", "/404.html", "the location of the 404 page, relative to the `static` directory.")
	flag.StringVar(&cssDirectory, "css", "/css", "the directory to serve css files from, relative to the `static` directory.")
	flag.StringVar(&javascriptDirectory, "js", "/js", "the directory to serve javascript files from, relative to the `static` directory.")
	flag.StringVar(&mediaDirectory, "media", "/media", "the directory to serve media files from, relative to the `static` directory.")
	flag.StringVar(&publicDirectory, "public", "/public", "the directory to serve public pages from, relative to the `static` directory.")
	flag.StringVar(&protectedDirectory, "protected", "/protected", "the directory to serve protected pages from, relative to the `static` directory.")

	flag.StringVar(&authClientId, "clientid", "", "the client ID for auth server.")
	flag.StringVar(&authClientSecret, "clientsecret", "", "the client secret for auth server.")
	flag.StringVar(&authClientIssuer, "issuer", "", "the issuer for auth server.")
	flag.StringVar(&authClientAudience, "audience", "", "the audience for auth server.")
	flag.StringVar(&callbackResourcePath, "callback", "/callback", "the callback for the auth server to use.")

	flag.StringVar(&fileStoreKey, "filestorekey", "", "the key to use for filestore encryption.")

	flag.StringVar(&serveHost, "host", "", "the host.")
	flag.StringVar(&servePort, "port", "443", "the port.")

	flag.StringVar(&bstApi, "api", "", "bst api host.")
	flag.StringVar(&bstApiBase, "apibase", "/", "bst api base path.")

	flag.Parse()
}