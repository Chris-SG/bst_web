package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func UserRouter() *mux.Router {
	userRouter := mux.NewRouter().PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("", UserProfile).Methods(http.MethodGet)

	return userRouter
}

func UserProfile(rw http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("./dist/user/user.html")
	rw.WriteHeader(200)
	rw.Write(fileBytes)
}