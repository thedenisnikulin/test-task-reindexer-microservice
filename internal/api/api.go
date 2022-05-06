package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reind01/internal"
	"reind01/internal/data"
	"reind01/internal/services"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
)

type Handler struct {
	Repo Repository[data.Author, int64]
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody internal.CreateAuthorRequest
	var authorModel data.Author
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

	err = h.Repo.Create(&authorModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	author, found := h.Repo.FindOne(int64(id))

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

	authors := h.Repo.GetAll(qty, page)

	var res internal.GetAllAuthorsResponse
	res.Authors = make([]*internal.GetAllAuthorsResponsePartial, len(authors))

	ch := make(chan *internal.GetAllAuthorsResponsePartial)

	for i := 0; i < len(res.Authors); i++ {
		go services.Process(authors[i], ch)
	}

	for i := 0; i < len(res.Authors); i++ {
		res.Authors[i] = <-ch
	}

	json.NewEncoder(w).Encode(res)
}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody internal.UpdateAuthorRequest
	var authorModel data.Author

	err := json.NewDecoder(r.Body).Decode(&authorReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	copier.Copy(&authorModel, &authorReqBody)

	err = h.Repo.Update(&authorModel)
	if err == data.NoItemsAffectedErr {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err == data.InternalServerErr {

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	err = h.Repo.Delete(int64(id))

	if err == data.NoItemsAffectedErr {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err == data.InternalServerErr {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
