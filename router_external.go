package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func CreateExternalRouters(middleware map[string]*negroni.Negroni) *mux.Router {
	externalRouter := mux.NewRouter().PathPrefix("/external").Subrouter()

	externalRouter.PathPrefix("/external/bst_api").Handler(negroni.New(
		negroni.Wrap(CreateBstApiRouter(middleware))))

	return externalRouter
}
