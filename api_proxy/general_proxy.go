package api_proxy

import (
	"bst_web/utilities"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chris-sg/bst_server_models"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
	"net/url"
)

// CreateBstApiRouter will generate a router mapped against BST API. Middleware
// may be passed in to then be used by certain routes.
func CreateBstApiRouter(prefix string, middleware map[string]*negroni.Negroni) *mux.Router {
	bstApiRouter := mux.NewRouter().PathPrefix(prefix + "/api").Subrouter()
	bstApiRouter.PathPrefix("/ddr").Handler(negroni.New(
		negroni.Wrap(CreateDdrProxy(prefix + "/api"))))
	bstApiRouter.PathPrefix("/drs").Handler(negroni.New(
		negroni.Wrap(CreateDrsProxy(prefix + "/api"))))
	bstApiRouter.Path("/status").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(StatusGet)))).Methods(http.MethodGet)
	bstApiRouter.Path("/bstuser").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(BstUserPut)))).Methods(http.MethodPut)
	bstApiRouter.Path("/eagate/login").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(EagateLoginGet)))).Methods(http.MethodGet)
	bstApiRouter.Path("/eagate/login").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(EagateLoginPost)))).Methods(http.MethodPost)
	bstApiRouter.Path("/eagate/logout").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(EagateLogoutPost)))).Methods(http.MethodPost)

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
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "status")

	status.Api = "bad"
	status.EaGate = "bad"
	status.Db = "bad"

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
	}
	res, err := utilities.GetClient().Do(req)
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

// StatusGet will call StatusGetImpl() and return the result.
func BstUserPut(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
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
	profile, err := utilities.ProfileForRequest(r)
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
	sub, ok := profile["sub"].(string)
	if !ok {
		status := bst_models.Status{
			Status:  "bad",
			Message: err.Error(),
		}

		bytes, _ := json.Marshal(status)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write(bytes)
		return
	}
	user := BstUserPutImpl(token, sub, r)

	bytes, _ := json.Marshal(user)
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
}

// StatusGetImpl will retrieve the current state of the api, the database and eagate.
func BstUserPutImpl(token string, sub string, r *http.Request) (userCache bst_models.UserCache) {
	utilities.ClearCacheValue("users", sub)
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "bstuser")

	req := &http.Request{
		Method:           http.MethodPut,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	request, err := ioutil.ReadAll(r.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(request))

	res, err := utilities.GetClient().Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &userCache)
	if err != nil {
		return
	}

	utilities.SetCacheValue("users", sub, userCache)
	return
}

func EagateLoginGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
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

	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "user/login")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := utilities.GetClient().Do(req)
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
	token, err := utilities.TokenForRequest(r)
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
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "user/login")

	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	b, _ := json.Marshal(loginRequest)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	res, err := utilities.GetClient().Do(req)
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
	token, err := utilities.TokenForRequest(r)
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
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "user/logout")

	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	b, _ := json.Marshal(logoutRequest)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	res, err := utilities.GetClient().Do(req)
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