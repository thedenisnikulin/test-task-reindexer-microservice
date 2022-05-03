package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reind01/internal/reindexerapp/models"
	"reind01/pkg/db"

	"github.com/jinzhu/copier"
)

type Handler struct {
	Db *db.Db
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody AuthorReqBody
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

	h.Db.Insert("authors", &authorModel)

	w.WriteHeader(200)
}

func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)

}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {

}
