package api

import (
	"github.com/gorilla/mux"
)

func SetRoutes(handler *Handler, r *mux.Router) {
	r.HandleFunc("/authors/{id}", handler.GetAuthor).
		Methods("GET").
		Headers("Content-Type", "application/json")

	r.HandleFunc("/authors/all/{qty}/{page}", handler.GetAllAuthors).
		Methods("GET").
		Headers("Content-Type", "application/json")

	r.HandleFunc("/authors", handler.CreateAuthor).
		Methods("POST")

	r.HandleFunc("/authors/{id}", handler.UpdateAuthor).
		Methods("PUT")

	r.HandleFunc("/authors/{id}", handler.DeleteAuthor).
		Methods("DELETE")
}
