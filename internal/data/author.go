package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"reind01/internal"
	"reind01/internal/infra"

	"github.com/coocood/freecache"
	"github.com/jinzhu/copier"
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
		LeftJoin(r.Db.Query(DbArticlesNamespaceName), "articles").
		On("articles_id", reindexer.SET, "id").
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
	it := r.Db.Query(DbAuthorsNamespaceName).
		Offset(int(qty*(page-1))).
		Limit(int(qty)).
		Sort("sort", true).
		Join(r.Db.Query(DbArticlesNamespaceName).Distinct("id"), "articles").
		On("articles_id", reindexer.SET, "id").
		Exec()

	models := make([]*Author, 0)

	for it.Next() {
		models = append(models, it.Object().(*Author))
	}

	return models
}

func (r *AuthorRepository) Create(reqModel *internal.CreateAuthorRequest) error {
	var authorModel Author
	var articleModels = make([]*Article, len(reqModel.Articles))

	err := copier.Copy(&authorModel, reqModel)
	if err != nil {
		r.Log.Error(err)
	}

	for i := 0; i < len(reqModel.Articles); i++ {
		var a Article
		copier.Copy(&a, reqModel.Articles[i])
		articleModels[i] = &a
	}

	// first insert an author without any articles to generate serial id
	insertedAuthors, err := r.Db.Insert(DbAuthorsNamespaceName,
		&authorModel,
		"id=serial()",
		"sort=serial()")

	if err != nil {
		r.Log.Error(err)
		return InternalServerErr
	}

	if insertedAuthors < 1 {
		return NoItemsAffectedErr
	}

	for i := 0; i < len(articleModels); i++ {
		articleModels[i].AuthorId = authorModel.Id // set generated serial author id
		insertedArticles, err := r.Db.Insert(DbArticlesNamespaceName, articleModels[i], "id=serial()")
		authorModel.Articles[i] = articleModels[i] // set updated articles id
		authorModel.ArticlesId = append(authorModel.ArticlesId, articleModels[i].Id)

		if err != nil {
			r.Log.Error(err)
			return InternalServerErr
		}

		if insertedArticles < 1 {
			return NoItemsAffectedErr
		}
	}

	_, err = r.Db.Update(DbAuthorsNamespaceName, authorModel) // then update articles
	if err != nil {
		return InternalServerErr
	}

	return nil
}

func (r *AuthorRepository) Update(reqModel *internal.UpdateAuthorRequest) error {
	deleted := r.Cache.Del([]byte(fmt.Sprintf("%v", reqModel.Id)))
	if !deleted {
		r.Log.Warnf("Cache couldn't be deleted, id=%v", reqModel.Id)
	}

	var author Author

	copier.Copy(&author, reqModel)
	for i := 0; i < len(reqModel.Articles); i++ {
		author.ArticlesId = append(author.ArticlesId, int64(reqModel.Articles[i].Id))
	}

	// update author
	updatedItems, err := r.Db.Update(DbAuthorsNamespaceName, &author)

	if err != nil {
		r.Log.Error(err)
		return InternalServerErr
	}

	if updatedItems < 1 {
		return NoItemsAffectedErr
	}

	// update articles
	for i := 0; i < len(reqModel.Articles); i++ {
		author.Articles[i].AuthorId = author.Id
		updatedItems, err = r.Db.Update(DbArticlesNamespaceName, author.Articles[i])

		if err != nil {
			r.Log.Error(err)
			return InternalServerErr
		}

		if updatedItems < 1 {
			return NoItemsAffectedErr
		}
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
		r.Log.Error(err)
		return InternalServerErr
	}

	if deletedItems < 1 {
		return NoItemsAffectedErr
	}

	_, err = r.Db.Query(DbArticlesNamespaceName).
		WhereInt64("author_id", reindexer.EQ, id).
		Delete()

	if err != nil {
		r.Log.Error(err)
		return InternalServerErr
	}

	return nil
}
