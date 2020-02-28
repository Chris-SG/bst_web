package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func AjaxRouter() *mux.Router {
	ajaxRouter := mux.NewRouter().PathPrefix("/ajax").Subrouter()

	ajaxRouter.HandleFunc("/apistatus", ApiStatus).Methods(http.MethodGet)

	return ajaxRouter
}

func ApiStatus(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte(StatusGet()))
}