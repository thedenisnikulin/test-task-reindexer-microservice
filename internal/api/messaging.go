package api

type CreateAuthorRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Articles []*struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"articles"`
}

type UpdateAuthorRequest struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Articles []*struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"articles"`
}

type GetAllAuthorsResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Articles []*struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	} `json:"articles"`
}
