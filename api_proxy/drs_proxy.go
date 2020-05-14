package api_proxy

import (
	"bst_web/utilities"
	"encoding/json"
	"fmt"
	bst_models "github.com/chris-sg/bst_server_models"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
	"net/url"
)

func CreateDrsProxy(prefix string) *mux.Router {
	drsProxy := mux.NewRouter().PathPrefix(prefix + "/drs").Subrouter()

	drsProxy.Path("/profile").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DrsProfilePatch)))).Methods(http.MethodPatch)

	drsProxy.Path("/details").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DrsDetailsGet)))).Methods(http.MethodGet)

	drsProxy.Path("/tabledata").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DrsTabledataGet)))).Methods(http.MethodGet)

	return drsProxy
}

func DrsProfilePatch(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	err = DrsProfilePatchImpl(token)

	bytes, _ := json.Marshal(err)
	if !err.Equals(bst_models.ErrorOK) {
		fmt.Printf("failed to update ddr profile: %s\n", err.Message)
	}

	rw.WriteHeader(err.CorrespondingHttpCode)
	rw.Write(bytes)
	return
}

func DrsProfilePatchImpl(token string) (err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "drs/profile")

	req := &http.Request{
		Method:           http.MethodPatch,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

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

func DrsDetailsGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	response, err := DrsDetailsGetImpl(token)

	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	rw.WriteHeader(err.CorrespondingHttpCode)
	rw.Write(response)
	return
}

func DrsDetailsGetImpl(token string) (response []byte, err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "drs/details")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	response, e = ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientResponse
		return
	}

	return
}

func DrsTabledataGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	response, err := DrsTabledataGetImpl(token)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
	return
}

func DrsTabledataGetImpl(token string) (response []byte, err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "drs/tabledata")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	response, e = ioutil.ReadAll(res.Body)
	if e != nil {
		err = bst_models.ErrorClientRequest
		return
	}

	if res.StatusCode != http.StatusOK {
		json.Unmarshal(response, &err)
		return
	}

	return
}