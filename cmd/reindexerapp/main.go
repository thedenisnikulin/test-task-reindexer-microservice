package main

import (
	"net/http"
	"reind01/config"
	"reind01/internal/api"
	"reind01/internal/data"
	database "reind01/internal/infra"

	"github.com/coocood/freecache"
	"github.com/gorilla/mux"
	"github.com/restream/reindexer"
	"github.com/sirupsen/logrus"
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

	db.OpenNamespace(
		data.DbAuthorsNamespaceName,
		reindexer.DefaultNamespaceOptions(),
		data.Author{})

	db.OpenNamespace(
		data.DbArticlesNamespaceName,
		reindexer.DefaultNamespaceOptions(),
		data.Article{})

	cache := freecache.NewCache(data.CacheSizeInBytes)


	logger := logrus.New()
	logger.SetReportCaller(true)

	repo := data.AuthorRepository{Db: &db, Cache: cache, Log: logger}
	handler := api.Handler{Repo: &repo}

	router := mux.NewRouter()
	api.SetRoutes(&handler, router)

	server := http.Server{
		Handler: router,
		Addr:    cfg.ServerConfig.Addr,
	}

	server.ListenAndServe()

	// TODO add signal/interruption handling
}
