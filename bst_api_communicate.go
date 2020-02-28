package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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

func Status_Get() string {
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

	type Status struct {
		Status string `json:"status"`
	}
	status := Status{}

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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OneTimePassword string `json:"otp,omitempty"`
}

func Eagate_Login_Post(token string, loginRequest LoginRequest) bool {
	data, err := json.Marshal(loginRequest)
	if err != nil {
		return false
	}
	uri, _ := url.Parse(bstApi + bstApiBase + "eagate/login")
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

func Ddr_Songs_Get() string {
	return ""
}