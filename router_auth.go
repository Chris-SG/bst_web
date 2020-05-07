package main

import (
	"bst_web/utilities"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

// AttachAuthRoutes attaches the router endpoints to the provided
// router that will be used in authentication.
func AttachAuthRoutes(r *mux.Router) {
	r.Path("/callback").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(http.HandlerFunc(utilities.CallbackHandler))))

	r.Path("/login").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(http.HandlerFunc(utilities.LoginHandler))))

	r.Path("/logout").Handler(utilities.GetCommonMiddleware().With(
		negroni.Wrap(http.HandlerFunc(utilities.LogoutHandler))))
}