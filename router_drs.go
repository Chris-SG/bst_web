package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func DrsRouter() *mux.Router {
	fmt.Println("Building drs routes...")
	ddrRouter := mux.NewRouter().PathPrefix("/drs").Subrouter()

	ddrRouter.HandleFunc("", DrsIndex).Methods(http.MethodGet)
	ddrRouter.HandleFunc("/stats", DrsStats).Methods(http.MethodGet)

	fmt.Println("Done")
	return ddrRouter
}

func DrsIndex(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/drs/drs.html")
	rw.WriteHeader(200)
	rw.Write(fileBytes)
}

func DrsStats(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/drs/stats.html")
	rw.WriteHeader(200)
	rw.Write(fileBytes)
}