package main

import (
	"bst_web/api_proxy"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func CreateExternalRouters(prefix string, middleware map[string]*negroni.Negroni) *mux.Router {
	externalRouter := mux.NewRouter().PathPrefix(prefix + "/external").Subrouter()
	externalRouter.PathPrefix("/api").Handler(negroni.New(
		negroni.Wrap(api_proxy.CreateBstApiRouter(prefix + "/external", middleware))))

	return externalRouter
}
