package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chris-sg/bst_server_models/bst_api_models"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/chris-sg/bst_server_models/bst_web_models"
)

var (
	client *http.Client
)

func InitClient() {
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:           nil,
		Timeout:       time.Second * 5,
	}
}

func StatusGet() string {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "status")
	fmt.Println(uri.String())
	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "bad"
	}

	status := bst_api_models.Status{}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(body)
	err = json.Unmarshal(body, &status)
	if err != nil {
		fmt.Println(err)
		return "unknown"
	}

	return status.Status
}

func EagateLoginGet(r *http.Request) {

}

func EagateLoginPost(token string, loginRequest bst_web_models.LoginRequest) bool {
	data, err := json.Marshal(loginRequest)
	if err != nil {
		return false
	}
	uri, _ := url.Parse(bstApi + bstApiBase + "user/login")
	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:           make(map[string][]string),
		Body:             ioutil.NopCloser(bytes.NewReader(data)),
		ContentLength:    int64(len(data)),
	}
	req.Header["Authorization"] = []string{"Bearer " + token}
	_, err = client.Do(req)

	return true
}

func DdrSongsGet() string {
	return ""
}