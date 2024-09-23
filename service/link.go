package service

import (
	"errors"
	"fmt"
	"markless/model"
	"markless/store"
	"markless/util"
	"net/url"
	"strings"
	"time"
)

func LinkCreate(user model.User, webURL string, desc string) (model.Link, error) {
	link := model.Link{}
	store.DB.Where("url = ?", webURL).Find(&link)
	if link.ID != 0 && link.UserID == user.ID {
		util.Logger.Warn("link重复:" + webURL)
		return link, errors.New("链接已存在")
	} else if link.ID != 0 && link.UserID != user.ID {
		existArch := model.Archive{}
		store.DB.Model(&link).Association("Archive").Find(&existArch)

		newLink := model.Link{}
		newLink.UserID = user.ID
		newLink.Url = webURL
		newLink.Title = link.Title
		newLink.Desc = desc
		newLink.Icon = link.Icon
		if existArch.Content != "" {
			newLink.Archive = &model.Archive{
				UpdateTime: time.Now(),
				Content:    existArch.Content,
			}
		}
		store.DB.Create(&newLink)
		return newLink, nil
	} else {
		newLink := model.Link{
			CreateTime: time.Now(),
			UserID:     user.ID,
			Url:        webURL,
		}
		pageInfo, err := util.ParsePage(webURL)
		if err != nil {
			return newLink, err
		} else {
			if pageInfo.Title != "" {
				newLink.Title = pageInfo.Title
			} else {
				newLink.Title = webURL
			}

			if desc != "" {
				newLink.Desc = desc
			} else if pageInfo.Desc != "" {
				newLink.Desc = pageInfo.Desc
			} else {
				newLink.Desc = webURL
			}

			if pageInfo.Icon != "" && strings.Contains(pageInfo.Icon, "http") {
				newLink.Icon = pageInfo.Icon
			} else {
				// 自行解析 URL拼接favicon
				parsedURL, err := url.Parse(webURL)
				if err != nil {
					util.Logger.Error(fmt.Println("Error parsing URL:", err))
				}
				rootPath := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
				newLink.Icon = rootPath + "/favicon.ico"
			}
			newLink.Archive = &model.Archive{
				UpdateTime: time.Now(),
				Content:    pageInfo.Content,
			}
			store.DB.Create(&newLink)
			return newLink, nil
		}
	}

}

func LinkAttachTag(user model.User, link *model.Link, tagNames []string) {
	tags := []model.Tag{}

	for _, v := range tagNames {
		if v == "" {
			continue
		}
		tag := model.Tag{Name: strings.Trim(v, " "), UserID: user.ID, CreateTime: time.Now()}
		store.DB.Where("name = ? AND user_id = ?", tag.Name, user.ID).Find(&tag)
		tags = append(tags, tag)
	}
	store.DB.Model(&link).Where("id = ? AND user_id =?", link.ID, user.ID).Association("Tags").Append(&tags)
}
