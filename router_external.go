package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func CreateExternalRouters(prefix string, middleware map[string]*negroni.Negroni) *mux.Router {
	externalRouter := mux.NewRouter().PathPrefix(prefix + "/external").Subrouter()
	externalRouter.PathPrefix("/bst_api").Handler(negroni.New(
		negroni.Wrap(CreateBstApiRouter(prefix + "/external", middleware))))

	return externalRouter
}
