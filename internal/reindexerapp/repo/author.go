package repo

import (
	"errors"
	"reind01/internal/reindexerapp"
	"reind01/internal/reindexerapp/models"

	"github.com/restream/reindexer"
)

type AuthorRepository struct {
	db *reindexer.Reindexer
}

func (r *AuthorRepository) Find(id int64) (model *models.Author, found bool) {
	author, found := r.db.Query(reindexerapp.DbAuthorsNamespaceName).
		WhereInt("id", reindexer.EQ, int(id)).
		Get()

	model = author.(*models.Author)

	return
}

func (r *AuthorRepository) Create(model *models.Author) error {
	inserted, err := r.db.Insert(reindexerapp.DbAuthorsNamespaceName, &model)

	if inserted != 1 {
		return errors.New("Item was not created.")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthorRepository) Update(model *models.Author) error {
	updated, err := r.db.Update(reindexerapp.DbAuthorsNamespaceName, model)

	if updated != 1 {
		return errors.New("Item was not updated.")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthorRepository) Delete(id int64) error {
	deleted, err := r.db.Query(reindexerapp.DbAuthorsNamespaceName).
		WhereInt64("id", reindexer.EQ, id).
		Delete()

	if deleted != 1 {
		return errors.New("Item was not deleted.")
	}

	if err != nil {
		return nil
	}

	return nil
}
