package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

func UserRouter() *mux.Router {
	userRouter := mux.NewRouter().PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("", UserProfile).Methods(http.MethodGet)

	return userRouter
}

func UserProfile(rw http.ResponseWriter, r *http.Request) {
	header := LoadHeader(r)
	footer := LoadFooter()

	session, _ := Store.Get(r, "auth-session")

	t, _:= template.ParseFiles("./dist/user_pages/user.html")
	replace := struct {
		Header string
		LoginForm string
		Footer string
	} {
		header,
		fmt.Sprint(session),
		footer,
	}
	fmt.Println("Working")
	rw.WriteHeader(200)
	t.Execute(rw, replace)
}