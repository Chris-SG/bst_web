package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"text/template"
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
	fileBytes, _ := ioutil.ReadFile("./dist/ddr_pages/ddr.html")
	fileText := string(fileBytes)

	t, _:= template.New("ddr").Parse(fileText)
	replace := struct {
		Header string
		Footer string
		CommonScripts string
		CommonSheets string
	} {
		LoadHeader(r),
		LoadFooter(),
		LoadCommonScripts(),
		LoadCommonSheets(),
	}

	rw.WriteHeader(200)
	t.Execute(rw, replace)
}

func DdrStats(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/ddr_pages/stats.html")
	fileText := string(fileBytes)

	t, _:= template.New("ddrstats").Parse(fileText)
	replace := struct {
		Header string
		Footer string
		CommonScripts string
		CommonSheets string
	} {
		LoadHeader(r),
		LoadFooter(),
		LoadCommonScripts(),
		LoadCommonSheets(),
	}

	rw.WriteHeader(200)
	t.Execute(rw, replace)
}