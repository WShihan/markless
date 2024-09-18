package api

import (
	"fmt"
	"markee/logging"
	"markee/model"
	"markee/store"
	"markee/tool"
	"markee/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func LinkAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	webURL := r.FormValue("url")
	tagNames := strings.Split(r.FormValue("tags"), "&")
	desc := r.FormValue("desc")
	if webURL == "" {
		logging.Logger.Info("link为空")
		tool.SetMsg(&w, "链接为空")
		util.Redirect(w, r, "/link/add")
		return
	}

	var title, icon string
	link := model.Link{}
	store.DB.Find(&link, "url = ?", webURL)

	newLink := model.Link{}
	// link重复
	if link.ID != 0 && link.UserID == user.ID {
		logging.Logger.Warn("link重复:" + webURL)
		tool.SetMsg(&w, "链接已存在")
		util.Redirect(w, r, "/")
		return
	} else if link.ID != 0 && link.UserID != user.ID {
		newLink.UserID = user.ID
		newLink.Url = webURL
		newLink.Title = link.Title
		newLink.Desc = link.Desc
		newLink.Icon = link.Icon

	} else {
		pageINfo, perr := util.Scrape(webURL, 10)
		if perr != nil {
			logging.Logger.Error(fmt.Println("Error scraping:", perr))
		} else {
			if pageINfo.Preview.Title == "" {
				title = webURL
			} else {
				title = pageINfo.Preview.Title
			}
			if desc == "" {
				if pageINfo.Preview.Description == "" {
					desc = r.FormValue("desc")
					if desc == "" {
						desc = title
					}
				} else {
					desc = pageINfo.Preview.Description
				}
			}

			if pageINfo.Preview.Icon == "" {
				// 解析 URL
				parsedURL, err := url.Parse(webURL)
				if err != nil {
					logging.Logger.Error(fmt.Println("Error parsing URL:", err))
				}
				// 尝试拼接favicon
				rootPath := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
				icon = rootPath + "/favicon.ico"
			} else {
				icon = pageINfo.Preview.Icon
			}
		}

		newLink.Title = title
		newLink.Desc = desc
		newLink.CreateTime = time.Now()
		newLink.Url = webURL
		newLink.Icon = icon
		newLink.UserID = user.ID
	}

	// 添加标签
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
	store.DB.Model(&newLink).Where("id = ? AND user_id =?", newLink.ID, user.ID).Association("Tags").Append(&tags)
	err := store.DB.Create(&newLink).Error
	if err != nil {
		tool.SetMsg(&w, "添加失败")
		logging.Logger.Error("添加失败" + err.Error())
		util.Redirect(w, r, "/link/add")
		return
	}
	logging.Logger.Info("添加书签成功：" + webURL)
	util.Redirect(w, r, "/")

}
func LinkUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, err := store.GetLinkByUser(user, int(id))
	if err != nil {
		tool.SetMsg(&w, "链接不存在")
	} else {
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
		err = store.DB.Save(&link).Error
		if err != nil {
			logging.Logger.Error("更新失败" + err.Error())
			tool.SetMsg(&w, "更新失败"+err.Error())
			return
		}
	}
	util.Redirect(w, r, "/")
}
func LinkRead(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = true
	err := store.DB.Save(&link).Error
	if err != nil {
		logging.Logger.Error("更新失败" + err.Error())
		model.ApiSuccess(&w, &model.ApiResponse{Msg: err.Error()})
		return
	}
	model.ApiSuccess(&w, &model.ApiResponse{Msg: "ok", Data: link})
}
func LinkUnread(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = false
	err := store.DB.Save(&link).Error
	if err != nil {
		logging.Logger.Error("更新失败" + err.Error())
		model.ApiSuccess(&w, &model.ApiResponse{Msg: err.Error()})
		return
	}
	model.ApiSuccess(&w, &model.ApiResponse{Msg: "ok", Data: link})
}

func LinkDel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	store.DB.Where("user_id = ?", user.ID).First(&link, id)
	err := store.DB.Unscoped().Delete(&link).Error
	if err != nil {
		logging.Logger.Error("删除失败" + err.Error())
		model.ApiSuccess(&w, &model.ApiResponse{Msg: err.Error()})
		return
	}
	logging.Logger.Info("删除书签成功：" + link.Url)
	model.ApiSuccess(&w, &model.ApiResponse{Msg: "ok", Data: link})
}
