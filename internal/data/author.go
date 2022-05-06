package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"reind01/internal/infra"

	"github.com/coocood/freecache"
	"github.com/restream/reindexer"
	"github.com/sirupsen/logrus"
)

var InternalServerErr = errors.New("Internal server error.")
var NoItemsAffectedErr = errors.New("No items affected error.")

type AuthorRepository struct {
	Db    *infra.Db
	Cache *freecache.Cache
	Log   *logrus.Logger
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

	author, found := r.Db.Query(DbAuthorsNamespaceName).
		WhereInt("id", reindexer.EQ, int(id)).
		Get()

	model := author.(*Author)

	marshaled, err := json.Marshal(model)

	err = r.Cache.Set(key, marshaled, CacheTtlInSecs)
	if err != nil {
		r.Log.Warnf("Cache couldn't be set, id=%v", id)
	}

	return model, found
}

func (r *AuthorRepository) GetAll(qty, page int) []*Author {
	// TODO use transactions
	it := r.Db.Query(DbAuthorsNamespaceName).
		Offset(int(qty*(page-1)+1)).
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
	insertedItems, err := r.Db.Insert(DbAuthorsNamespaceName, model,
		"id=serial()",
		"articles.id=serial()",
		"sort=serial()")

	if err != nil {
		return InternalServerErr
	}

	if insertedItems < 1 {
		return NoItemsAffectedErr
	}

	return nil
}

func (r *AuthorRepository) Update(model *Author) error {
	deleted := r.Cache.Del([]byte(fmt.Sprintf("%v", model.Id)))
	if !deleted {
		r.Log.Warnf("Cache couldn't be deleted, id=%v", model.Id)
	}

	updatedItems, err := r.Db.Update(DbAuthorsNamespaceName, model)

	if updatedItems < 1 {
		return NoItemsAffectedErr
	}

	if err != nil {
		return InternalServerErr
	}

	return nil
}

func (r *AuthorRepository) Delete(id int64) error {
	affected := r.Cache.Del([]byte(fmt.Sprintf("%v", id)))

	if !affected {
		r.Log.Warnf("Cache couldn't be deleted, id=%v", id)
	}

	deletedItems, err := r.Db.Query(DbAuthorsNamespaceName).
		WhereInt64("id", reindexer.EQ, id).
		Delete()

	if err != nil {
		return InternalServerErr
	}

	if deletedItems < 1 {
		return NoItemsAffectedErr
	}

	return nil
}
