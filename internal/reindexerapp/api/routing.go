package api

import (
	"github.com/gorilla/mux"
)

func SetRoutes(handler *Handler, r *mux.Router) {
	r.HandleFunc("/authors/{id}", handler.GetAuthor).Methods("GET")
	r.HandleFunc("/author/all/{qty}/{page}", handler.GetAllAuthors).Methods("GET")
	r.HandleFunc("/authors", handler.CreateAuthor).Methods("POST")
	r.HandleFunc("/authors/{id}", handler.UpdateAuthor).Methods("PUT")
	r.HandleFunc("/authors/{id}", handler.DeleteAuthor).Methods("DELETE")
}
