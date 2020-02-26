package main

import (
	"bytes"
	"fmt"
	"text/template"
	"net/http"
)

type HeaderTemplate struct {
	LoginState string
	Ops []struct {
		Link string
		Text string
	}
}

func LoadHeader(r *http.Request) (string, error) {
	t, err := template.ParseFiles("./templates/header.html")
	if err != nil {
		panic(err)
	}
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		panic(err)
	}

	headerTemplate := HeaderTemplate{}
	var header bytes.Buffer

	if session != nil {
		if _, ok := session.Values["access_token"]; ok {
			headerTemplate.LoginState = `<button onclick="location.href='/logout';">Logout</button>`

			t.Execute(&header, headerTemplate)
			return header.String(), nil
		}
	}
	headerTemplate.LoginState = `<button onclick="location.href='/login';">Login</button>`

	t.Execute(&header, headerTemplate)
	fmt.Println(header.String())
	return header.String(), nil
}