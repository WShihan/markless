package injection

import "markee/model"

type Search struct {
	Keyword    string
	PrePage    int
	Page       int
	Limit      int
	NextPage   int
	TagName    string
	ReadStatus int // 0：所有 1:已读 2:未读
	Count      int
}

type LinkEditInjection struct {
	Link model.Link
	Tags []model.Tag
}

type LinkPage struct {
	Title string
	Env   Env
	Search
	Data       interface{}
	TagStastic map[string]int
}
