package utilities

import (
	"net/http"
	"time"
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
		Timeout:       time.Second * 60,
	}
}

func GetClient() *http.Client {
	return bstApiClient
}