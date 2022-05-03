package api

type CreateAuthorReqBody struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Articles []*struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"articles"`
}

type UpdateAuthorReqBody struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Articles []*struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"articles"`
}
