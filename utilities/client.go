package utilities

import (
	"net/http"
)

var (
	bstApiClient *http.Client
)

func InitClient() {
	bstApiClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:           nil,
		Timeout:       0,
	}
}

func GetClient() *http.Client {
	return bstApiClient
}