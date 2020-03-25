package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func DdrRouter() *mux.Router {
	fmt.Println("Building ddr routes...")
	ddrRouter := mux.NewRouter().PathPrefix("/ddr").Subrouter()

	ddrRouter.HandleFunc("", DdrIndex).Methods(http.MethodGet)
	ddrRouter.HandleFunc("/stats", DdrStats).Methods(http.MethodGet)

	fmt.Println("Done")
	return ddrRouter
}

func DdrIndex(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/ddr/ddr.html")
	rw.WriteHeader(200)
	rw.Write(fileBytes)
}

func DdrStats(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/ddr/stats.html")
	rw.WriteHeader(200)
	rw.Write(fileBytes)
}