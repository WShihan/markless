package injection

import (
	"markless/model"
	"time"
)

type TagStatstic struct {
	ID         int
	Name       string
	Count      int
	CreateTime time.Time
	Links      []model.Link
}

type Search struct {
	Keyword    string
	PrePage    int
	Page       int
	Pages      int
	Limit      int
	NextPage   int
	TagName    string
	ReadStatus int // 0：所有 1:已阅 2:未阅
	Count      int
}

type LinkEditInjection struct {
	Link model.Link
	Tags []model.Tag
}

type LinkPage struct {
	Env  Env
	Page PageInjection
	Search
	Data       interface{}
	TagStastic map[string]*TagStatstic
}
