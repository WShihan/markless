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
	Keyword    string `json:"keyword"`
	PrePage    int    `json:"pre_page"`
	Page       int    `json:"page"`
	Pages      int    `json:"pages"`
	Limit      int    `json:"limit"`
	NextPage   int    `json:"next_page"`
	TagName    string `json:"tag"`
	ReadStatus int    `json:"read_state"` // 0：所有 1:已阅 2:未阅
	Count      int    `json:"count"`
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
