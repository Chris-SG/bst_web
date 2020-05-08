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

	return drsProxy
}

func DrsProfilePatch(rw http.ResponseWriter, r *http.Request) {
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

	status := DrsProfilePatchImpl(token)

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

func DrsProfilePatchImpl(token string) (status bst_models.Status) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "drs/profile")

	req := &http.Request{
		Method:           http.MethodPatch,
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
	json.Unmarshal(body, &status)

	return
}

func DrsDetailsGet(rw http.ResponseWriter, r *http.Request) {
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

	response, code := DrsDetailsGetImpl(token)

	rw.WriteHeader(code)
	rw.Write(response)
	return
}

func DrsDetailsGetImpl(token string) (response []byte, code int) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "drs/details")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := utilities.GetClient().Do(req)
	code = res.StatusCode
	if err != nil {
		response = []byte(`{"status":"bad","message":"api_err"}`)
		return
	}

	defer res.Body.Close()
	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		response = []byte(`{"status":"bad","message":"api_err"}`)
		return
	}

	return
}