package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reind01/internal"
	"reind01/internal/data"
	"reind01/internal/services"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Handler struct {
	Repo AuthorRepository
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody internal.CreateAuthorRequest
	err := json.NewDecoder(r.Body).Decode(&authorReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.Create(&authorReqBody)
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
	wg := new(sync.WaitGroup)

	for i := 0; i < len(res.Authors); i++ {
		wg.Add(1)
		go services.Process(authors[i], ch, wg)
	}

	quit := make(chan bool)
	go func() {
		wg.Wait()
		quit <- true
	}()

	for {
		select {
		case out := <-ch:
			res.Authors = append(res.Authors, out)
		case <-quit:
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	// for i := 0; i < len(res.Authors); i++ {
	// 	res.Authors[i] = <-ch
	// }

}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	var authorReqBody internal.UpdateAuthorRequest

	err := json.NewDecoder(r.Body).Decode(&authorReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.Update(&authorReqBody)
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
