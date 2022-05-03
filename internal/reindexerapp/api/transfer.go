package api

type AuthorReqBody struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Articles []*struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"articles"`
}
