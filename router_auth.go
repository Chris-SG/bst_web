package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

func AttachAuthRoutes(r *mux.Router) {
	r.Path("/callback").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(CallbackHandler))))

	r.Path("/login").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(LoginHandler))))

	r.Path("/logout").Handler(commonMiddleware.With(
		negroni.Wrap(http.HandlerFunc(LogoutHandler))))
}