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

	token, _ := TokenForRequest(r)
	status, users := EagateLoginGetImpl(token)

	t, _:= template.New("user").Parse(fileText)
	replace := struct {
		Header string
		LoginForm string
		Footer string
		CommonScripts string
		CommonSheets string
		LoggedIn bool
		EagateUsername string
	} {
		LoadHeader(r),
		fmt.Sprint(session),
		LoadFooter(),
		LoadCommonScripts(),
		LoadCommonSheets(),
		false,
		"",
	}
	if status.Status == "ok" && len(users) > 0 {
		replace.LoggedIn = !users[0].Expired
		replace.EagateUsername = users[0].Username
	}

	rw.WriteHeader(200)
	t.Execute(rw, replace)
}