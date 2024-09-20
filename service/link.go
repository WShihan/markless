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
	store.DB.Where("url = ? AND user_id = ?", webURL, user.ID).Find(&link)
	if link.ID != 0 && link.UserID == user.ID {
		util.Logger.Warn("link重复:" + webURL)
		return link, errors.New("链接已存在")
	} else if link.ID != 0 && link.UserID != user.ID {
		newLink := model.Link{}
		newLink.UserID = user.ID
		newLink.Url = webURL
		newLink.Title = link.Title
		newLink.Desc = desc
		newLink.Icon = link.Icon
		return newLink, nil
	} else {
		newLink := model.Link{
			CreateTime: time.Now(),
			UserID:     user.ID,
			Url:        webURL,
		}
		pageINfo, perr := util.Scrape(webURL, 10)
		if perr != nil {
			return newLink, perr
		} else {
			if pageINfo.Preview.Title != "" {
				newLink.Title = pageINfo.Preview.Title
			} else {
				newLink.Title = webURL
			}

			if desc != "" {
				newLink.Desc = desc
			} else if pageINfo.Preview.Description != "" {
				newLink.Desc = pageINfo.Preview.Description
			} else {
				newLink.Desc = webURL
			}

			if pageINfo.Preview.Icon != "" && strings.Contains(pageINfo.Preview.Icon, "http") {
				newLink.Icon = pageINfo.Preview.Icon
			} else {
				// 自行解析 URL拼接favicon
				parsedURL, err := url.Parse(webURL)
				if err != nil {
					util.Logger.Error(fmt.Println("Error parsing URL:", err))
				}
				rootPath := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
				newLink.Icon = rootPath + "/favicon.ico"
			}
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
		if tag.ID != 0 {
			util.Logger.Info("该标签当前用户已存在：" + tag.Name)
			continue
		}
		tags = append(tags, tag)
	}
	store.DB.Model(&link).Where("id = ? AND user_id =?", link.ID, user.ID).Association("Tags").Append(&tags)
}
