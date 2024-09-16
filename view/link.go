package view

import (
	"fmt"
	"markee/model"
	"markee/store"
	"markee/util"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func LinkAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	webURL := r.FormValue("url")
	tagNames := strings.Split(r.FormValue("tags"), "&")
	pageINfo, _ := util.Scrape(webURL, 10)
	var title, icon, desc string
	if pageINfo.Preview.Title == "" {
		title = "未知标题"
	} else {
		title = pageINfo.Preview.Title
	}
	if pageINfo.Preview.Description == "" {
		desc = r.FormValue("desc")
		if desc == "" {
			desc = title
		}
	} else {
		desc = pageINfo.Preview.Description
	}
	if pageINfo.Preview.Icon == "" {
		// 解析 URL
		parsedURL, err := url.Parse(webURL)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}
		// 尝试拼接favicon
		rootPath := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
		icon = rootPath + "/favicon.ico"
	} else {
		icon = pageINfo.Preview.Icon
	}

	link := model.Link{}
	user := model.User{}

	store.DB.First(&user)
	store.DB.Find(&link, "url = ?", webURL)
	if link.Url == "" {
		link.Title = title
		link.Desc = desc
		link.CreatedAt = time.Now()
		link.Url = webURL
		link.Icon = icon
		link.UserID = user.ID
		tags := []model.Tag{}
		for _, v := range tagNames {
			if v == "" {
				continue
			}
			tag := model.Tag{Name: strings.Trim(v, " "), UserID: user.ID, CreateTime: time.Now()}
			store.DB.Find(&tag, "name = ?", v)
			if tag.ID == 0 {
				store.DB.Create(&tag)
			}
			tags = append(tags, tag)
		}
		store.DB.Model(&link).Association("Tags").Append(&tags)
		store.DB.Create(&link)
		Redirect(w, r, "/")

	} else {
		Redirect(w, r, "/link/add")
	}

}
func LinkUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	link := model.Link{}
	store.DB.First(&link, id)
	link.Title = r.FormValue("title")
	link.Desc = r.FormValue("desc")
	link.Url = r.FormValue("url")
	tagArr := strings.Split(r.FormValue("tags"), "&")
	tags := []model.Tag{}
	for _, v := range tagArr {
		if v == "" {
			continue
		}
		tag := model.Tag{}
		store.DB.Find(&tag, "name = ?", strings.Trim(v, " "))
		tags = append(tags, tag)
	}
	store.DB.Model(&link).Association("Tags").Append(&tags)
	store.DB.Save(&link)
	Redirect(w, r, "/")
}
func LinkRead(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	link := model.Link{}
	store.DB.First(&link, id)
	link.Read = !link.Read
	store.DB.Save(&link)
	Redirect(w, r, "/")
}

func LinkDel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	link := model.Link{}
	store.DB.First(&link, id)
	store.DB.Delete(&link)
	Redirect(w, r, "/")
}
