package main

import (
	"net/http"
	"reind01/internal/reindexerapp/api"
	"reind01/internal/reindexerapp/models"
	"reind01/pkg/config"
	"reind01/pkg/db"

	"github.com/restream/reindexer"
	"github.com/gorilla/mux"
)

func main() {
	var cfg *config.Config
	dotEnvPath := ".env"
	configYamlPath := "config.yml"

	cfg, err := config.NewConfigYaml(configYamlPath)
	if err != nil {
		cfg, err = config.NewConfigDotEnv(dotEnvPath)
		if err != nil {
			cfg, err = config.NewConfigEnv()
			if err != nil {
				panic("No configuration is provided")
			}
		}
	}


	database := db.OpenDb(&cfg.DbConfig)

	if err := database.Ping(); err != nil {
		panic(err)
	}

	found, err := database.HasNamespace("authors")
	if err != nil {
		panic(err)
	}
	if !found {
		database.OpenNamespace("authors", reindexer.DefaultNamespaceOptions(), models.Author{})
	}


	router := mux.NewRouter()
	api.SetRoutes(router)

	server := http.Server{
		Handler: router,
		Addr:    cfg.ServerConfig.Addr,
	}

	server.ListenAndServe()
}
