package store

import "markee/model"

func TagStat() map[string]int {
	user := model.User{}
	DB.First(&user)
	tags := []model.Tag{}
	links := []model.Link{}
	DB.Model(&user).Association("Links").Find(&links)
	DB.Model(&user).Association("Tags").Find(&tags)
	staMap := make(map[string]int)
	// 添加标签
	for _, v := range tags {
		staMap[v.Name] = 0
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
			staMap[vv.Name]++
		}
	}

	return staMap
}
