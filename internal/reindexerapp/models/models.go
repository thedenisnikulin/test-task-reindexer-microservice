package models

type Author struct {
	Id       int64     `reindex:"id,,pk"`
	Name     string    `reindex:"name"`
	Age      int       `reindex:"age"`
	Articles []Article `reindex:"articles"`
}

type Article struct {
	Id    int64  `reindex:"id,,pk"`
	Title string `reindex:"title"`
	Body  string `reindex:"body"`
}
