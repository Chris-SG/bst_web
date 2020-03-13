package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"text/template"
)

func DdrRouter() *mux.Router {
	userRouter := mux.NewRouter().PathPrefix("/ddr").Subrouter()

	userRouter.HandleFunc("", DdrIndex).Methods(http.MethodGet)

	return userRouter
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