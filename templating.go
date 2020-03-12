package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"
)

type ClickBox struct {
	Link string
	Text string
	Class string
}

type HeaderTemplate struct {
	HomePage string
	DropdownText string
	Ops []ClickBox
}

func LoadHeader(r *http.Request) string {
	t, err := template.ParseFiles("./templates/header.html")
	if err != nil {
		panic(err)
	}
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		panic(err)
	}

	headerTemplate := HeaderTemplate{}
	headerTemplate.HomePage = "https://" + serveHost
	var header bytes.Buffer

	if session != nil {
		if _, ok := session.Values["access_token"]; ok {
			var nickname string
			profileMap, ok := session.Values["profile"].(map[string]interface{})
			if !ok {
				nickname = "No Profile"
			} else {
				nickname, ok = profileMap["nickname"].(string)
				if !ok {
					nickname = "No Nickname"
				}
			}
			headerTemplate.DropdownText = nickname
			headerTemplate.Ops = append(headerTemplate.Ops, ClickBox{Link: "/user", Text: "Profile", Class: "btn-profile"})
			headerTemplate.Ops = append(headerTemplate.Ops, ClickBox{Link: "/logout", Text: "Logout", Class: "btn-logout"})

			t.Execute(&header, headerTemplate)
			return header.String()
		}
	}
	headerTemplate.DropdownText = "Not logged in."
	headerTemplate.Ops = append(headerTemplate.Ops, ClickBox{Link: "/login", Text: "Login", Class: "btn-login"})

	t.Execute(&header, headerTemplate)
	return header.String()
}
type ContentTemplate struct {
	Footer string
}

func LoadFooter() string {

	f, _ := ioutil.ReadFile("./templates/footer.html")
	return string(f)
}

func LoadCommonScripts() string {
	return `<script src='https://code.jquery.com/jquery-latest.min.js' type='text/javascript'></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
<script src='/js/external/js.cookie.js'></script>
<script src='/js/common.js'></script>`
}

func LoadCommonSheets() string {
	return `
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <link rel="stylesheet" href="/css/common.css">`
}

func LoadEagateLogin(user string) {

}