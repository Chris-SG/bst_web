package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func DDRSongs(rw http.ResponseWriter, r *http.Request) {

}

func DDRSongsId(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
}