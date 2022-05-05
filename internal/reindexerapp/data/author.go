package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"reind01/internal/reindexerapp"
	"reind01/pkg/db"

	"github.com/coocood/freecache"
	"github.com/restream/reindexer"
)

type AuthorRepository struct {
	Db    *db.Db
	Cache *freecache.Cache
}

func (r *AuthorRepository) FindOne(id int64) (*Author, bool) {
	key := []byte(fmt.Sprintf("%v", id))
	bytes, err := r.Cache.Get(key)
	if err == nil {
		var fromCache Author
		err := json.Unmarshal(bytes, &fromCache)
		if err == nil {
			return &fromCache, true
		}
	}

	author, found := r.Db.Query(reindexerapp.DbAuthorsNamespaceName).
		WhereInt("id", reindexer.EQ, int(id)).
		Get()

	model := author.(*Author)

	marshaled, err := json.Marshal(model)

	err = r.Cache.Set(key, marshaled, reindexerapp.CacheTtlInSecs)
	if err != nil {
		// TODO log cache set failed
	}

	return model, found
}

func (r *AuthorRepository) GetAll(qty, page int) []*Author {
	it := r.Db.Query(reindexerapp.DbAuthorsNamespaceName).
		Offset(int(qty*(page-1) + 1)).
		Limit(int(qty)).
		Sort("sort", true).
		Exec()

	models := make([]*Author, 0)

	for it.Next() {
		models = append(models, it.Object().(*Author))
	}

	return models
}

func (r *AuthorRepository) Create(model *Author) error {
	insertedItems, err := r.Db.Insert(reindexerapp.DbAuthorsNamespaceName, &model)

	if insertedItems != 1 {
		return errors.New("Item was not created.")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthorRepository) Update(model *Author) error {
	deleted := r.Cache.Del([]byte(fmt.Sprintf("%v", model.Id)))
	if !deleted {
		// TODO log
	}

	updatedItems, err := r.Db.Update(reindexerapp.DbAuthorsNamespaceName, model)

	if updatedItems != 1 {
		return errors.New("Item was not updated.")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthorRepository) Delete(id int64) error {
	affected := r.Cache.Del([]byte(fmt.Sprintf("%v", id)))

	if !affected {
		// TODO log
	}

	deletedItems, err := r.Db.Query(reindexerapp.DbAuthorsNamespaceName).
		WhereInt64("id", reindexer.EQ, id).
		Delete()


	// TODO create aliases for the following errors & match them in http handler
	// to make http status code more comprehensive
	if deletedItems != 1 {
		return errors.New("Item was not deleted.")
	}

	if err != nil {
		return err
	}

	return nil
}
