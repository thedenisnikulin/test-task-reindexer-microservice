package api

import "reind01/internal/data"
import "reind01/internal"

type AuthorRepository interface {
	FindOne(id int64) (*data.Author, bool)
	GetAll(qty, page int) []*data.Author
	Create(model *internal.CreateAuthorRequest) error
	Update(model *internal.UpdateAuthorRequest) error
	Delete(id int64) error
}
