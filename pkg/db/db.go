package db

import (
	"fmt"
	"github.com/restream/reindexer"
	"reind01/pkg/config"
	//_ "github.com/restream/reindexer/bindings/builtin"
)

type Db struct {
	*reindexer.Reindexer
}

func OpenDb(config *config.DbConfig) Db {
	reidx := reindexer.NewReindex(
		//fmt.Sprintf("cproto://%v:%v@%v/%v", config.DbUser, config.DbPass, config.DbAddr, config.DbName),
		fmt.Sprintf("cproto://%v/%v", config.DbAddr, config.DbName), reindexer.WithCreateDBIfMissing())

	return Db{reidx}
}


func (db *Db) HasNamespace(ns string) (bool, error) {
	namespaces, err := db.Reindexer.DescribeNamespaces()
	if err != nil {
		return false, err
	}

	for i := 0; i < len(namespaces); i++ {
		if namespaces[i].Name == ns {
			return true, nil
		}
	}

	return false, nil
}
