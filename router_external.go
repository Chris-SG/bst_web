package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func CreateExternalRouters(middleware map[string]*negroni.Negroni) *mux.Router {
	externalRouter := mux.NewRouter().PathPrefix("/external").Subrouter()
	externalRouter.PathPrefix("/external/bst_api").Handler(negroni.New(
		negroni.Wrap(CreateBstApiRouter(middleware))))

	return externalRouter
}

func walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	pathTemplate, err := route.GetPathTemplate()
	if err == nil {
		fmt.Println("ROUTE:", pathTemplate)
	}
	pathRegexp, err := route.GetPathRegexp()
	if err == nil {
		fmt.Println("Path regexp:", pathRegexp)
	}
	for _, r2 := range ancestors {
		pathTemplate, err = r2.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err = r2.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
	}
	return nil
}