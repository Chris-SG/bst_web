package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"text/template"
)

func UserRouter() *mux.Router {
	userRouter := mux.NewRouter().PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("", UserProfile).Methods(http.MethodGet)

	return userRouter
}

func UserProfile(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/user_pages/user.html")
	fileText := string(fileBytes)

	session, _ := Store.Get(r, "auth-session")

	t, _:= template.New("user").Parse(fileText)
	replace := struct {
		Header string
		LoginForm string
		Footer string
		CommonScripts string
		CommonSheets string
	} {
		LoadHeader(r),
		fmt.Sprint(session),
		LoadFooter(),
		LoadCommonScripts(),
		LoadCommonSheets(),
	}

	rw.WriteHeader(200)
	t.Execute(rw, replace)
}