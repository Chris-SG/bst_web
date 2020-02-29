package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func UserRouter() *mux.Router {
	ajaxRouter := mux.NewRouter().PathPrefix("/user").Subrouter()

	ajaxRouter.HandleFunc("/", UserProfile).Methods(http.MethodGet)

	return ajaxRouter
}

func UserProfile(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte(StatusGet()))
}