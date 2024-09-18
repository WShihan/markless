package store

import "markee/model"

func GetLinkByUser(user model.User, id int) (model.Link, error) {
	link := model.Link{}
	err := DB.Where("user_id = ? AND id = ?", user.ID, id).Find(&link).Error
	return link, err
}

func LinkExist(url string, user model.User) bool {
	link := model.Link{}
	err := DB.Where("user_id = ?", user.ID).Order("id desc").Find(&link).Error
	if err != nil || link.Url != "" {
		return false
	} else {
		return true
	}
}
