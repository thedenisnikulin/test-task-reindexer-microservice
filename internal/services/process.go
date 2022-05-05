package services

import (
	"reind01/internal"
	"reind01/internal/data"

	"github.com/jinzhu/copier"
)

func Process(model *data.Author, ch chan *internal.GetAllAuthorsResponse) {
	var r internal.GetAllAuthorsResponse
	copier.Copy(model, &r)
	ch <- &r
}
