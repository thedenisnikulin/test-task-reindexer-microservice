package data

// 2 level nestedness (collection of authors where each author contains a collection of articles)

type Author struct {
	Id         int64      `reindex:"id,,pk"`
	Name       string     `reindex:"name"`
	Age        int        `reindex:"age"`
	ArticlesId []int64    `reindex:"articles_id"`
	Articles   []*Article `reindex:"articles,,joined"`
	Sort       int        `reindex:"sort,tree"`
}

type Article struct {
	Id       int64  `reindex:"id,,pk"`
	Title    string `reindex:"title"`
	Body     string `reindex:"body,text"`
	AuthorId int64  `reindex:"author_id"`
}
