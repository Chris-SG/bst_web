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
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/profile/update")

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

func DdrRefreshPatch(rw http.ResponseWriter, r *http.Request) {
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
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/profile/refresh")

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

func DdrStatsGet(rw http.ResponseWriter, r *http.Request) {
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

	stats := DdrStatsGetImpl(token)

	rw.Write([]byte(stats))
	return
}

func DdrStatsGetImpl(token string) (stats string) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/songs/scores/extended")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := utilities.GetClient().Do(req)
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

func DdrProfileGet(rw http.ResponseWriter, r *http.Request) {
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

	stats := DdrProfileGetImpl(token)
	if len(stats) == 0 {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(stats))
	return
}

func DdrProfileGetImpl(token string) (profile string) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/profile")

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := utilities.GetClient().Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	profile = string(body)
	return
}

func DdrSongScoresGet(rw http.ResponseWriter, r *http.Request) {
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

	buf := make([]byte, 0)
	r.Body.Read(buf)

	response, status := DdrSongScoresGetImpl(token, r.URL.RawQuery)

	rw.WriteHeader(status)
	rw.Write([]byte(response))
	return
}

func DdrSongScoresGetImpl(token string, queryParams string) (response string, code int) {
	uri, _ := url.Parse("https://" + utilities.BstApi + utilities.BstApiBase + "ddr/song/scores")
	uri.RawQuery = queryParams

	req := &http.Request{
		Method:           http.MethodGet,
		URL:              uri,
		Header:			  make(map[string][]string),
	}
	req.Header.Add("Authorization", "Bearer " + token)

	res, err := utilities.GetClient().Do(req)
	code = res.StatusCode
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	response = string(body)
	return
}