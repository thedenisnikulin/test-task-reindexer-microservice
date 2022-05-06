package services

import (
	"reind01/internal"
	"reind01/internal/data"
	"sync"

	"github.com/jinzhu/copier"
)

func Process(model *data.Author, ch chan<- *internal.GetAllAuthorsResponsePartial, wg *sync.WaitGroup) {
	defer wg.Done()
	var r internal.GetAllAuthorsResponsePartial
	copier.Copy(&r, model)
	ch <- &r
}
