package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reind01/internal/reindexerapp"
	"reind01/internal/reindexerapp/models"
	"reind01/pkg/db"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/restream/reindexer"
)

type Handler struct {
	Db *db.Db
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody CreateAuthorReqBody
	var authorModel models.Author
	err := json.NewDecoder(r.Body).Decode(&authorReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = copier.Copy(&authorModel, &authorReqBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.Db.Insert(reindexerapp.DbAuthorsNamespaceName, &authorModel)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	author, found := h.Db.Query(reindexerapp.DbAuthorsNamespaceName).
		WhereInt("id", reindexer.EQ, int(id)).
		Get()
	
	author = author.(*models.Author)

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(author)
}

func (h *Handler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qty, err1 := strconv.Atoi(vars["qty"])
	page, err2 := strconv.Atoi(vars["page"])
	if err1 != nil || err2 != nil {
		http.Error(w, fmt.Sprintf("%v; %v", err1, err2), http.StatusBadRequest)
		return
	}

	it := h.Db.Query(reindexerapp.DbAuthorsNamespaceName).
		Offset(int(qty * (page - 1) + 1)).
		Limit(int(qty)).
		Exec()
	
	authors := make([]*models.Author, 0)

	for it.Next() {
		authors = append(authors, it.Object().(*models.Author))
	}

	json.NewEncoder(w).Encode(authors)
}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody UpdateAuthorReqBody
	var authorModel models.Author

	err := json.NewDecoder(r.Body).Decode(&authorReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	copier.Copy(&authorModel, &authorReqBody)

	updated, err := h.Db.Update(reindexerapp.DbAuthorsNamespaceName, authorModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if updated == 0 {
		http.Error(w, "No item was updated.", http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleted, err := h.Db.Query(reindexerapp.DbAuthorsNamespaceName).
		WhereInt("id", reindexer.EQ, id).
		Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		http.Error(w, "No item was deleted", http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusOK)
}
