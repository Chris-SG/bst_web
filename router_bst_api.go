package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chris-sg/bst_server_models"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
	"net/url"
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
		Timeout:       time.Second * 45,
	}
}

// CreateBstApiRouter will generate a router mapped against BST API. Middleware
// may be passed in to then be used by certain routes.
func CreateBstApiRouter(prefix string, middleware map[string]*negroni.Negroni) *mux.Router {
	bstApiRouter := mux.NewRouter().PathPrefix(prefix + "/bst_api").Subrouter()
	bstApiRouter.Path("/status").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(StatusGet)))).Methods(http.MethodGet)
	bstApiRouter.Path("/eagate_login").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(EagateLoginPost)))).Methods(http.MethodPost)
	bstApiRouter.Path("/eagate_logout").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(EagateLogoutPost)))).Methods(http.MethodPost)

	bstApiRouter.Path("/ddr_update").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrUpdatePatch)))).Methods(http.MethodPatch)
	bstApiRouter.Path("/ddr_refresh").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrRefreshPatch)))).Methods(http.MethodPatch)
	bstApiRouter.Path("/ddr_stats").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrStatsGet)))).Methods(http.MethodGet)


	return bstApiRouter
}

// StatusGet will call StatusGetImpl() and return the result.
func StatusGet(rw http.ResponseWriter, r *http.Request) {
	status := StatusGetImpl()

	bytes, _ := json.Marshal(status)
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
}

// StatusGetImpl will retrieve the current state of the api, the database and eagate.
func StatusGetImpl() (status bst_models.ApiStatus) {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "status")

	status.Api = "bad"
	status.EaGate = "bad"
	status.Db = "bad"

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
	}
	res, err := bstApiClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &status)
	if err != nil {
		status.Api = "unknown"
	}

	return
}

func EagateLoginGet(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}
	status, users := EagateLoginGetImpl(token)

	if status.Status == "bad" {
		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(bytes)
		return
	}

	bytes, _ := json.Marshal(users)
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	return
}

func EagateLoginGetImpl(token string) (status bst_models.Status, users []bst_models.EagateUser){

	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "user/login")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := bstApiClient.Do(req)
	if err != nil {
		status.Status = "bad"
		status.Message = "api error"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	users = make([]bst_models.EagateUser, 0)
	json.Unmarshal(body, &users)

	status.Status = "ok"
	status.Message = fmt.Sprintf("found %d users", len(users))
	return
}

func EagateLoginPost(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	fmt.Println(string(body))

	loginRequest := bst_models.LoginRequest{}
	json.Unmarshal(body, &loginRequest)

	status := EagateLoginPostImpl(token, loginRequest)

	bytes, _ := json.Marshal(status)
	if status.Status == "ok" {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.Write(bytes)
	return
}

func EagateLoginPostImpl(token string, loginRequest bst_models.LoginRequest) (status bst_models.Status) {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "user/login")

	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	b, _ := json.Marshal(loginRequest)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	res, err := bstApiClient.Do(req)
	if err != nil {
		status.Status = "bad"
		status.Message = "api error"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &status)

	return
}

func EagateLogoutPost(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	logoutRequest := bst_models.LogoutRequest{}
	err = json.Unmarshal(body, &logoutRequest)
	fmt.Println(err)
	fmt.Printf("%s\n", body)
	fmt.Println(logoutRequest)

	status := EagateLogoutPostImpl(token, logoutRequest)

	bytes, _ := json.Marshal(status)
	if status.Status == "ok" {
		rw.WriteHeader(http.StatusOK)
	} else {
		fmt.Printf("failed to logout user: %s\n", status.Message)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.Write(bytes)
	return
}

func EagateLogoutPostImpl(token string, logoutRequest bst_models.LogoutRequest) (status bst_models.Status) {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "user/logout")

	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	b, _ := json.Marshal(logoutRequest)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	res, err := bstApiClient.Do(req)
	if err != nil {
		status.Status = "bad"
		status.Message = "api error"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &status)

	return
}

func DdrUpdatePatch(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	status := DdrUpdatePatchImpl(token)

	bytes, _ := json.Marshal(status)
	if status.Status == "ok" {
		rw.WriteHeader(http.StatusOK)
	} else {
		fmt.Printf("failed to update ddr profile: %s\n", status.Message)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.Write(bytes)
	return
}

func DdrUpdatePatchImpl(token string) (status bst_models.Status) {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "ddr/profile/update")

	req := &http.Request{
		Method:           http.MethodPatch,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := bstApiClient.Do(req)
	if err != nil {
		status.Status = "bad"
		status.Message = "api error"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &status)

	return
}

func DdrRefreshPatch(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	status := DdrRefreshPatchImpl(token)

	bytes, _ := json.Marshal(status)
	if status.Status == "ok" {
		rw.WriteHeader(http.StatusOK)
	} else {
		fmt.Printf("failed to refresh ddr profile: %s\n", status.Message)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.Write(bytes)
	return
}

func DdrRefreshPatchImpl(token string) (status bst_models.Status) {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "ddr/profile/refresh")

	req := &http.Request{
		Method:           http.MethodPatch,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := bstApiClient.Do(req)
	if err != nil {
		status.Status = "bad"
		status.Message = "api error"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &status)

	return
}

func DdrStatsGet(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}

	stats := DdrStatsGetImpl(token)

	rw.Write([]byte(stats))
	return
}

func DdrStatsGetImpl(token string) (stats string) {
	uri, _ := url.Parse("https://" + bstApi + bstApiBase + "ddr/songs/scores/extended")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := bstApiClient.Do(req)
	if err != nil {
		stats = "<a>API Error</a>"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		stats = "<a>API Error</a>"
		return
	}
	stats = string(body)
	return
}