package api_proxy

import (
	"bst_web/utilities"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chris-sg/bst_server_models"
	"github.com/golang/glog"
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

func AddImpersonateToRequest(r *http.Request, req *http.Request) {
	if c, err := r.Cookie("impersonate"); err == nil && len(c.Value) > 0 {
		glog.Infof("%s impersonate cookie is set", c.String())
		req.Header.Set("Impersonate-User", c.Value)
	}
}


// StatusGet will call StatusGetImpl() and return the result.
func StatusGet(rw http.ResponseWriter, r *http.Request) {
	status, err := StatusGetImpl()
	if !err.Equals(bst_models.ErrorOK) {
		glog.Error(err)
	}

	bytes, _ := json.Marshal(status)
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
}

// StatusGetImpl will retrieve the current state of the api, the database and eagate.
func StatusGetImpl() (status bst_models.ApiStatus, err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "status")

	status.Api = "bad"
	status.EaGate = "bad"
	status.Db = "bad"

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
	}
	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientResponse
		return
	}

	e = json.Unmarshal(body, &status)
	if e != nil {
		err = bst_models.ErrorJsonDecode
		return
	}

	return
}

// StatusGet will call StatusGetImpl() and return the result.
func BstUserPut(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	profile, err := utilities.ProfileForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	sub, ok := profile["sub"].(string)
	if !ok {
		bytes, _ := json.Marshal(bst_models.ErrorJwtProfile)
		rw.WriteHeader(bst_models.ErrorJwtProfile.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}
	user, err := BstUserPutImpl(token, sub, r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	bytes, _ := json.Marshal(user)
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	return
}

// StatusGetImpl will retrieve the current state of the api, the database and eagate.
func BstUserPutImpl(token string, sub string, r *http.Request) (userCache bst_models.UserCache, err bst_models.Error) {
	err = bst_models.ErrorOK
	utilities.ClearCacheValue("users", sub)
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "bstuser")

	req := &http.Request{
		Method:           http.MethodPut,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	request, e := ioutil.ReadAll(r.Body)
	if e != nil {
		err = bst_models.ErrorBadBody
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(request))

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientResponse
		return
	}

	e = json.Unmarshal(body, &userCache)
	if e != nil {
		err = bst_models.ErrorJsonDecode
		return
	}

	utilities.SetCacheValue("users", sub, userCache)
	return
}

func EagateLoginGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}
	err, users := EagateLoginGetImpl(token, r)

	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	bytes, _ := json.Marshal(users)
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
	return
}

func EagateLoginGetImpl(token string, r *http.Request) (err bst_models.Error, users []bst_models.EagateUser){
	err = bst_models.ErrorOK

	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "user/login")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)
	AddImpersonateToRequest(r, req)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientResponse
		return
	}

	users = make([]bst_models.EagateUser, 0)
	json.Unmarshal(body, &users)

	return
}

func EagateLoginPost(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	defer r.Body.Close()
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		bytes, _ := json.Marshal(bst_models.ErrorBadBody)
		rw.WriteHeader(bst_models.ErrorBadBody.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	loginRequest := bst_models.LoginRequest{}
	json.Unmarshal(body, &loginRequest)

	err = EagateLoginPostImpl(token, loginRequest)

	bytes, _ := json.Marshal(err)
	if !err.Equals(bst_models.ErrorOK) {}

	rw.WriteHeader(err.CorrespondingHttpCode)
	rw.Write(bytes)
	return
}

// TODO: use form instead of body
func EagateLoginPostImpl(token string, loginRequest bst_models.LoginRequest) (err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "user/login")

	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	b, _ := json.Marshal(loginRequest)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientResponse
		return
	}
	json.Unmarshal(body, &err)

	return
}

// TODO: use form instead of body
func EagateLogoutPost(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	defer r.Body.Close()
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		bytes, _ := json.Marshal(bst_models.ErrorBadRequest)
		rw.WriteHeader(bst_models.ErrorBadRequest.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	logoutRequest := bst_models.LogoutRequest{}
	e = json.Unmarshal(body, &logoutRequest)
	if e != nil {
		bytes, _ := json.Marshal(bst_models.ErrorJsonDecode)
		rw.WriteHeader(bst_models.ErrorJsonDecode.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	err = EagateLogoutPostImpl(token, logoutRequest)
	b, e := json.Marshal(err)
	if e != nil {
		b, _ := json.Marshal(bst_models.ErrorJsonEncode)
		rw.WriteHeader(bst_models.ErrorJsonEncode.CorrespondingHttpCode)
		rw.Write(b)
		return
	}

	if !err.Equals(bst_models.ErrorOK) {
		fmt.Printf("failed to logout user: %s\n", err.Message)
	}

	rw.WriteHeader(err.CorrespondingHttpCode)
	rw.Write(b)
	return
}

func EagateLogoutPostImpl(token string, logoutRequest bst_models.LogoutRequest) (err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "user/logout")

	req := &http.Request{
		Method:           http.MethodPost,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	b, _ := json.Marshal(logoutRequest)
	req.Body = ioutil.NopCloser(bytes.NewReader(b))

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientResponse
		return
	}
	json.Unmarshal(body, &err)

	return
}