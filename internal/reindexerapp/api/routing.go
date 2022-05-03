package api

import (
	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) {
	r.HandleFunc("/author/{id}", GetAuthor).Methods("GET")
	r.HandleFunc("/author", CreateAuthor).Methods("POST")
	r.HandleFunc("/author", UpdateAuthor).Methods("PUT")
	r.HandleFunc("/author", DeleteAuthor).Methods("DELETE")
}
