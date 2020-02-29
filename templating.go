package main

import (
	"bytes"
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

type FooterTemplate struct {
	BstApiStatus string
}

func LoadFooter() string {
	t, err := template.ParseFiles("./templates/footer.html")
	if err != nil {
		panic(err)
	}

	var footer bytes.Buffer

	footerTemplate := FooterTemplate{
		BstApiStatus: "unknown",
	}

	t.Execute(&footer, footerTemplate)
	return footer.String()
}

func LoadEagateLogin(user string) {
}