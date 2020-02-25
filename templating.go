package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

func LoadHeader(r *http.Request) (string, error) {
	t, err := template.ParseFiles("./templates/header.html")
	if err != nil {
		return "", err
	}
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		return "", err
	}

	replacement := make(map[string]interface{})
	replacement["Ops"] = make([]string, 0)

	if _, ok := session.Values["access_token"]; ok {
		var writer io.Writer
		replacement["LoginState"] = `<button onclick="location.href='/logout';">Logout</button>`

		t.Execute(writer, replacement)
		var outText []byte
		writer.Write(outText)
		fmt.Println(outText)
		return string(outText), nil
	}

	var writer io.Writer
	replacement["LoginState"] = `<button onclick="location.href='/login';">Login</button>`

	t.Execute(writer, replacement)
	var result string
	io.WriteString(writer, result)
	fmt.Println(result)
	return result, nil
}