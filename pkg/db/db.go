package db

import (
	"fmt"
	"reind01/internal/reindexerapp/models"
	"reind01/pkg/config"

	"github.com/restream/reindexer"
	_ "github.com/restream/reindexer/bindings/builtin"
)

type Db struct {
	*reindexer.Reindexer
}

func Init(config *config.Config) Db {
	reidx := reindexer.NewReindex(
		fmt.Sprintf("cproto://%v:%v@%v/%v", config.DbUser, config.DbPass, config.DbAddr, config.DbName),
		reindexer.WithCreateDBIfMissing())

	reidx.OpenNamespace("authors", reindexer.DefaultNamespaceOptions(), models.Author{})

	return Db { reidx }
}
