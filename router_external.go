package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
)

func CreateExternalRouters(middleware map[string]*negroni.Negroni) *mux.Router {
	externalRouter := mux.NewRouter().PathPrefix("/external").Subrouter()

	externalRouter.PathPrefix("/bst_api").Handler(negroni.New(
		negroni.Wrap(CreateBstApiRouter(middleware))))

	externalRouter.Path("/status").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(StatusGet))))


	fmt.Println("-----WALKING ROUTES-----")
	externalRouter.Walk(walk)
	fmt.Println("-----FINISHED WALK-----")

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
	if !strings.HasSuffix(pathRegexp, "$") {
		route.Subrouter().Walk(walk)
	}
	fmt.Println()
	return nil
}