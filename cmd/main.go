package main

import (
	"net/http"
	"reind01/internal/reindexerapp/api"
	"reind01/internal/reindexerapp/models"
	"reind01/pkg/config"
	database "reind01/pkg/db"

	"github.com/gorilla/mux"
	"github.com/restream/reindexer"
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

	db := database.OpenDb(&cfg.DbConfig)
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	found, err := db.HasNamespace("authors")
	if err != nil {
		panic(err)
	}
	if !found {
		db.OpenNamespace("authors", reindexer.DefaultNamespaceOptions(), models.Author{})
	}

	handler := api.Handler{Db: &db}

	router := mux.NewRouter()
	api.SetRoutes(&handler, router)

	server := http.Server{
		Handler: router,
		Addr:    cfg.ServerConfig.Addr,
	}

	server.ListenAndServe()
}
