package services

import (
	"reind01/internal"
	"reind01/internal/data"

	"github.com/jinzhu/copier"
)

func Process(model *data.Author, ch chan<- *internal.GetAllAuthorsResponsePartial) {
	var r internal.GetAllAuthorsResponsePartial
	copier.Copy(&r, model)
	ch <- &r
}
