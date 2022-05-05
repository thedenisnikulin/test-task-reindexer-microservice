package data

// 2 level nestedness (collection of authors where each author contains a collection of articles)

type Author struct {
	Id       int64      `reindex:"id,,pk"`
	Name     string     `reindex:"name"`
	Age      int        `reindex:"age"`
	Articles []*Article `reindex:"articles"`
	Sort     int        `reindex:"sort,tree"`
}

type Article struct {
	Id    int64  `reindex:"id"`
	Title string `reindex:"title"`
	Body  string `reindex:"body,text"`
}
