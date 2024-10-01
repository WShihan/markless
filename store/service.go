package store

import (
	"markless/injection"
	"markless/model"
)

func TagStat(user model.User) map[string]*injection.TagStatstic {
	tags := []model.Tag{}
	links := []model.Link{}
	DB.Model(&user).Association("Links").Find(&links)
	DB.Model(&user).Association("Tags").Find(&tags)
	staMap := make(map[string]*injection.TagStatstic)
	// 添加标签
	for _, v := range tags {
		staMap[v.Name] = &injection.TagStatstic{
			ID:         v.ID,
			Name:       v.Name,
			Count:      0,
			CreateTime: v.CreateTime,
		}
	}

	// 获取links标签
	for i, v := range links {
		tags := []model.Tag{}
		DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}

	// 统计标签
	for _, v := range links {
		for _, vv := range v.Tags {
			staMap[vv.Name].Count += 1
		}
	}

	return staMap
}
