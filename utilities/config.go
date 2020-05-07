package utilities

import "flag"

var (
	StaticDirectory string
	IndexPage string
	NotFoundPage string
	CssDirectory string
	JavascriptDirectory string
	MediaDirectory string
	PublicDirectory string
	ProtectedDirectory string

	authClientId string
	authClientSecret string
	authClientIssuer string
	authClientAudience string
	callbackResourcePath string

	fileStoreKey string

	ServeHost string
	ServePort string

	BstApi string
	BstApiBase string
)

// LoadConfig populates general configuration values to be used with the program.
func LoadConfig() {
	flag.StringVar(&StaticDirectory, "static", "./dist", "the directory containing all static files.")
	flag.StringVar(&IndexPage, "index", "/index.html", "the location of the index page, relative to the `static` directory.")
	flag.StringVar(&NotFoundPage, "404", "/404.html", "the location of the 404 page, relative to the `static` directory.")
	flag.StringVar(&CssDirectory, "css", "/css", "the directory to serve css files from, relative to the `static` directory.")
	flag.StringVar(&JavascriptDirectory, "js", "/js", "the directory to serve javascript files from, relative to the `static` directory.")
	flag.StringVar(&MediaDirectory, "media", "/media", "the directory to serve media files from, relative to the `static` directory.")
	flag.StringVar(&PublicDirectory, "public", "/public", "the directory to serve public pages from, relative to the `static` directory.")
	flag.StringVar(&ProtectedDirectory, "protected", "/protected", "the directory to serve protected pages from, relative to the `static` directory.")

	flag.StringVar(&authClientId, "clientid", "", "the client ID for auth server.")
	flag.StringVar(&authClientSecret, "clientsecret", "", "the client secret for auth server.")
	flag.StringVar(&authClientIssuer, "issuer", "", "the issuer for auth server.")
	flag.StringVar(&authClientAudience, "audience", "", "the audience for auth server.")
	flag.StringVar(&callbackResourcePath, "callback", "/callback", "the callback for the auth server to use.")

	flag.StringVar(&fileStoreKey, "filestorekey", "", "the key to use for filestore encryption.")

	flag.StringVar(&ServeHost, "host", "", "the host.")
	flag.StringVar(&ServePort, "port", "443", "the port.")

	flag.StringVar(&BstApi, "api", "", "bst api host.")
	flag.StringVar(&BstApiBase, "apibase", "/", "bst api base path.")

	flag.Parse()
}