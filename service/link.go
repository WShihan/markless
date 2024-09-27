package service

import (
	"errors"
	"fmt"
	"markless/injection"
	"markless/model"
	"markless/store"
	"markless/util"
	"math"
	"net/url"
	"strings"
	"time"
)

var (
	LIMIT = 20
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
func detectPages(count int, limit int) int {
	res := float64(count) / float64(limit)
	return int(math.Ceil(res))
}

func FilterLinkByTag(opt *injection.Search, user model.User) []model.Link {
	rawLinks := []model.Link{}
	offset := opt.Page * opt.Limit
	targetTag := model.Tag{}
	store.DB.Where("name = ? and user_id", opt.TagName, user.ID).Find(&targetTag)
	store.DB.Model(&targetTag).Where("user_id = ?", user.ID).Association("Links").Find(&rawLinks)
	links := rawLinks[:0]
	if opt.ReadStatus == 0 {
		links = append(links, rawLinks...)
	} else if opt.ReadStatus == 1 {
		for _, v := range rawLinks {
			if v.Read {
				links = append(links, v)
			}
		}
	} else if opt.ReadStatus == 2 {
		for _, v := range rawLinks {
			if !v.Read {
				links = append(links, v)
			}
		}
	}
	opt.Count = len(links)
	opt.Pages = detectPages(len(links), opt.Limit)
	opt.Keyword = "#" + opt.TagName

	if int(offset) < len(links) {
		end := int(offset) + opt.Limit
		if end > len(links) {
			end = len(links)
		}
		links = links[offset:end]
	}

	// 绑定标签
	for i, v := range links {
		tags := []model.Tag{}
		store.DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}

	return links
}

func FlterLinkByKeyword(opt *injection.Search, user model.User) []model.Link {
	links := []model.Link{}
	offset := opt.Page * opt.Limit
	var count int64
	condition := "%" + opt.Keyword + "%"
	var err error
	if opt.ReadStatus == 0 {
		store.DB.Model(&model.Link{}).Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Count(&count)
		err = store.DB.Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Order("created_at DESC").Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
	} else if opt.ReadStatus == 1 {
		store.DB.Model(&model.Link{}).Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", true).Count(&count)
		err = store.DB.Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", true).Order("created_at DESC").Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
	} else {
		store.DB.Model(&model.Link{}).Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", false).Count(&count)
		err = store.DB.Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", false).Limit(opt.Limit).Order("created_at DESC").Offset(int(offset)).Find(&links).Error
	}
	if err != nil {
		util.Logger.Error(err.Error())
	}
	opt.Count = int(count)
	opt.Pages = detectPages(opt.Count, opt.Limit)

	// 绑定标签
	for i, v := range links {
		tags := []model.Tag{}
		store.DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}

	return links
}
