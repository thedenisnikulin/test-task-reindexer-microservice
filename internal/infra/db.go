package infra

import (
	"fmt"
	"github.com/restream/reindexer"
	"reind01/config"
)

type Db struct {
	*reindexer.Reindexer
}

func OpenDb(config *config.DbConfig) Db {
	reidx := reindexer.NewReindex(
		//fmt.Sprintf("cproto://%v:%v@%v/%v", config.DbUser, config.DbPass, config.DbAddr, config.DbName),
		// TODO handle user/pass
		fmt.Sprintf("cproto://%v/%v", config.DbAddr, config.DbName), reindexer.WithCreateDBIfMissing())

	return Db{reidx}
}
