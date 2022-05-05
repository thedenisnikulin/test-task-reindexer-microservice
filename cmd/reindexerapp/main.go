package main

import (
	"net/http"
	"reind01/internal"
	"reind01/internal/api"
	"reind01/internal/data"
	"reind01/config"
	database "reind01/internal/infra"

	"github.com/coocood/freecache"
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

	db.OpenNamespace(
		internal.DbAuthorsNamespaceName,
		reindexer.DefaultNamespaceOptions(),
		data.Author{})

	cache := freecache.NewCache(internal.CacheSizeInBytes)

	// for i := int64(0); i < 50; i++ {
	// 	db.Insert(internal.DbAuthorsNamespaceName, &models.Author{
	// 		Id: i,
	// 		Name: fmt.Sprintf("name #%v", i),
	// 		Age: int(i) + rand.Int(),
	// 		Articles: []*models.Article{{
	// 			Id: i + rand.Int63(),
	// 			Title: fmt.Sprintf("damn #%v", i),
	// 			Body: fmt.Sprintf("Some body %v", i),
	// 		}},
	// 	})
	// }

	repo := data.AuthorRepository{Db: &db, Cache: cache}
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
