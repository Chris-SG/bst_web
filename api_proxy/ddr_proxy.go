package api_proxy

import (
	"bst_web/utilities"
	"encoding/json"
	"fmt"
	bst_models "github.com/chris-sg/bst_server_models"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
	"net/url"
)

func CreateDdrProxy(prefix string) *mux.Router {
	ddrProxy := mux.NewRouter().PathPrefix(prefix + "/ddr").Subrouter()

	ddrProxy.Path("/profile/update").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrUpdatePatch)))).Methods(http.MethodPatch)
	ddrProxy.Path("/profile/refresh").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrRefreshPatch)))).Methods(http.MethodPatch)
	ddrProxy.Path("/stats").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrStatsGet)))).Methods(http.MethodGet)
	ddrProxy.Path("/profile").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrProfileGet)))).Methods(http.MethodGet)
	ddrProxy.Path("/song/scores").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(DdrSongScoresGet)))).Methods(http.MethodGet)

	return ddrProxy
}



func DdrUpdatePatch(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	err = DdrUpdatePatchImpl(token, r)

	bytes, _ := json.Marshal(err)
	if !err.Equals(bst_models.ErrorOK) {
		fmt.Printf("failed to update ddr profile: %s\n", err.Message)
	}

	rw.WriteHeader(err.CorrespondingHttpCode)
	rw.Write(bytes)
	return
}

func DdrUpdatePatchImpl(token string, r *http.Request) (err bst_models.Error) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/profile/update")

	req := &http.Request{
		Method:           http.MethodPatch,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)
	AddImpersonateToRequest(r, req)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorApiInaccessible
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorApiInaccessible
		return
	}
	json.Unmarshal(body, &err)

	return
}

func DdrRefreshPatch(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	err = DdrRefreshPatchImpl(token, r)

	bytes, _ := json.Marshal(err)
	if !err.Equals(bst_models.ErrorOK) {
		fmt.Printf("failed to refresh ddr profile: %s\n", err.Message)
	}

	rw.WriteHeader(err.CorrespondingHttpCode)
	rw.Write(bytes)
	return
}

func DdrRefreshPatchImpl(token string, r *http.Request) (err bst_models.Error) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/profile/refresh")

	req := &http.Request{
		Method:           http.MethodPatch,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)
	AddImpersonateToRequest(r, req)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientResponse
		return
	}
	json.Unmarshal(body, &err)

	return
}

func DdrStatsGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	stats, err := DdrStatsGetImpl(token, r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	rw.Write([]byte(stats))
	rw.WriteHeader(http.StatusOK)
	return
}

func DdrStatsGetImpl(token string, r *http.Request) (stats string, err bst_models.Error) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/songs/scores/extended")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)
	AddImpersonateToRequest(r, req)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientResponse
		return
	}

	if res.StatusCode != http.StatusOK {
		json.Unmarshal(body, &err)
		return
	}

	stats = string(body)
	return
}

func DdrProfileGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	stats, err := DdrProfileGetImpl(token, r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}
	if len(stats) == 0 {
		bytes, _ := json.Marshal(bst_models.ErrorDdrStats)
		rw.WriteHeader(bst_models.ErrorDdrStats.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(stats))
	return
}

func DdrProfileGetImpl(token string, r *http.Request) (profile string, err bst_models.Error) {
	err = bst_models.ErrorOK
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/profile")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)
	AddImpersonateToRequest(r, req)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientResponse
		return
	}

	if res.StatusCode != http.StatusOK {
		json.Unmarshal(body, &err)
		return
	}

	profile = string(body)
	return
}

// TODO: query string instead of body
func DdrSongScoresGet(rw http.ResponseWriter, r *http.Request) {
	token, err := utilities.TokenForRequest(r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	buf := make([]byte, 0)
	r.Body.Read(buf)

	response, err := DdrSongScoresGetImpl(token, r.URL.RawQuery, r)
	if !err.Equals(bst_models.ErrorOK) {
		bytes, _ := json.Marshal(err)
		rw.WriteHeader(err.CorrespondingHttpCode)
		rw.Write(bytes)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(response))
	return
}

func DdrSongScoresGetImpl(token string, queryParams string, r *http.Request) (response string, err bst_models.Error) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/song/scores")
	uri.RawQuery = queryParams

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)
	AddImpersonateToRequest(r, req)

	res, e := utilities.GetClient().Do(req)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientRequest
		return
	}

	defer res.Body.Close()
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		glog.Errorf("failed request: %s", e.Error())
		err = bst_models.ErrorClientResponse
		return
	}

	if res.StatusCode != http.StatusOK {
		json.Unmarshal(body, &err)
		return
	}

	response = string(body)
	return
}
