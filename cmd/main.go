package main

import (
	"net/http"
	"reind01/internal/reindexerapp/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	api.SetRoutes(router)

	server := http.Server {
		Handler: router,
		Addr: "127.0.0.1:8000", // TODO from config
	}

	server.ListenAndServe()
}
