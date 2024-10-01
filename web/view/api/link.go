package api

import (
	"encoding/json"
	"markless/injection"
	"markless/local"
	"markless/model"
	"markless/service"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/server"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	LIMIT = 20
)

type LinkAddPost struct {
	Url  string `json:"url"`
	Desc string `json:"desc"`
	Tags string `json:"tags"`
	Read bool   `json:"read"`
}

type LinkAllData struct {
	Links  []model.Link     `json:"links"`
	Search injection.Search `json:"search"`
}

type LinkMarkLinkPost struct {
	Links []int `json:"links"`
	Read  bool  `json:"read"`
}

type LinkExistPost struct {
	Url string `json:"url"`
}

type LinkUpdatePost struct {
	ID    int    `json:"id"`
	Url   string `json:"url"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Read  bool   `json:"read"`
}

type LinkAttachTagsPost struct {
	ID   int      `json:"id"`
	Tags []string `json:"tags"`
}

func LinkAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	err := store.DB.Preload("Links").Find(&user).Error
	if err != nil {
		server.ApiFailed(&w, 200, err.Error())
		return
	}
	server.ApiSuccess(&w, &user.Links)
}
func LinkAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	postBody := LinkAddPost{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postBody)
	if err != nil {
		panic(err)
	}
	if postBody.Url == "" {
		panic(local.Translate("tip.link.empty", user.Lang))
	}
	link, err := service.LinkCreate(user, postBody.Url, postBody.Desc)
	if err != nil {
		panic(err)
	}
	link.Read = postBody.Read
	service.LinkAttachTag(user, &link, strings.Split(postBody.Tags, "&"))

	util.Logger.Info("add link success" + postBody.Url)
	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok", Data: link})

}

func LinkOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id := params.ByName("id")
	link := model.Link{}
	store.DB.Preload("Archive").Where("id = ? AND user_id =?", id, user.ID).Find(&link)
	store.DB.Preload("Tags").Where("id = ? AND user_id =?", id, user.ID).Find(&link)

	server.ApiSuccess(&w, link)

}

func LinkPagination(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")
	readVal := r.URL.Query().Get("read")
	tagName := ""
	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	read, _ := strconv.ParseInt(readVal, 10, 64)
	pagenum = pagenum - 1

	var links []model.Link
	if strings.Contains(keyword, "#") {
		tagName = strings.Replace(keyword, "#", "", 1)
	}
	searchOpt := injection.Search{
		PrePage:    int(pagenum) - 1,
		Page:       int(pagenum),
		Limit:      LIMIT,
		NextPage:   int(pagenum) + 2,
		TagName:    tagName,
		Keyword:    keyword,
		ReadStatus: int(read),
	}
	if tagName != "" {
		links = service.FilterLinkByTag(&searchOpt, user)
	} else {
		links = service.FlterLinkByKeyword(&searchOpt, user)
	}
	searchOpt.Page++
	data := LinkAllData{
		Links:  links,
		Search: searchOpt,
	}
	server.ApiSuccess(&w, &data)
}

func LinkRead(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = true
	err := store.DB.Save(&link).Error
	if err != nil {
		server.ApiSuccess(&w, &server.ApiResponse{Msg: err.Error()})
		return
	}
	server.ApiSuccess(&w, &link)
}

func LinkUnread(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = false
	err := store.DB.Save(&link).Error
	if err != nil {
		util.Logger.Error("update link failed" + err.Error())
		server.ApiSuccess(&w, nil)
		return
	}
	server.ApiSuccess(&w, &link)
}

func LinkDel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	store.DB.Where("user_id = ?", user.ID).First(&link, id)
	// err := store.DB.Unscoped().Delete(&link).Error
	err := store.DB.Select("Archive").Unscoped().Delete(&link).Error
	if err != nil {
		util.Logger.Error("delete link failed" + err.Error())
		server.ApiSuccess(&w, &server.ApiResponse{Msg: err.Error()})
		return
	}
	util.Logger.Info("delete link success" + link.Url)
	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok", Data: link})
}

func LinkExist(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	postData := LinkExistPost{}
	err := tool.ConvertJSON2Struct(&postData, r)
	if postData.Url == "" || err != nil {
		server.ApiFailed(&w, 200, "Link does not exist!")
		return
	}
	url := postData.Url
	link := model.Link{}
	store.DB.Where("url = ? AND user_id = ?", url, user.ID).Find(&link)
	if link.Url != "" {
		server.ApiSuccess(&w, nil)
		return
	}
	server.ApiFailed(&w, 200, "Link does not exist!")
}

func MarkAllAsReadOrRead(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	post := LinkMarkLinkPost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		server.ApiFailed(&w, 200, err.Error())
		return
	}
	store.DB.Model(&model.Link{}).Where("user_id = ? and id in ?", user.ID, post.Links).Update("read", post.Read)
	server.ApiSuccess(&w, nil)
}

func LinkUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	post := LinkUpdatePost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		server.ApiFailed(&w, 201, "参数错误")
		return
	}
	link := model.Link{}
	err = store.DB.Where("id =? AND user_id=?", post.ID, user.ID).First(&link).Error
	if err != nil {
		server.ApiFailed(&w, 201, "参数"+err.Error())
		return
	}
	link.Url = post.Url
	link.Title = post.Title
	link.Desc = post.Desc
	link.Read = post.Read
	err = store.DB.Save(&link).Error
	if err != nil {
		server.ApiFailed(&w, 201, "参数"+err.Error())
		return
	}
	server.ApiSuccess(&w, nil)

}

func LinkAttachTags(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	link := model.Link{}
	post := LinkAttachTagsPost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		server.ApiFailed(&w, 201, "参数错误")
		return
	}

	err = store.DB.Where("id=? AND user_id=?", post.ID, user.ID).First(&link).Error
	if err != nil {
		server.ApiFailed(&w, 201, "错误"+err.Error())
		return
	}

	updatedTags := []model.Tag{}
	for _, v := range post.Tags {
		if v == "" {
			continue
		}
		tag := model.Tag{
			Name:       v,
			UserID:     user.ID,
			CreateTime: time.Now(),
		}
		store.DB.Where("name=? AND user_id=?", v, user.ID).First(&tag)
		updatedTags = append(updatedTags, tag)
	}
	store.DB.Model(&link).Association("Tags").Clear()
	err = store.DB.Model(&link).Association("Tags").Append(&updatedTags)
	if err != nil {
		server.ApiFailed(&w, 201, err.Error())
	}
	server.ApiSuccess(&w, nil)
}

func LinkUpdateArchive(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id := params.ByName("id")
	link := model.Link{}
	store.DB.Preload("Archive").Where("user_id = ? AND id = ?", user.ID, id).First(&link)
	if link.Url == "" {
		server.ApiFailed(&w, 201, "书签不存在")
		return
	}
	pageInfo, err := util.ParsePage(link.Url)
	if err != nil {
		server.ApiFailed(&w, 201, "解析书签异常"+err.Error())
		return
	}
	if link.Archive != nil {
		link.Archive.Content = pageInfo.Content
		link.Archive.UpdateTime = time.Now()
	} else {
		link.Archive = &model.Archive{
			LinkID:     link.ID,
			Content:    pageInfo.Content,
			UpdateTime: time.Now(),
		}
	}

	err = store.DB.Save(&link.Archive).Error
	if err != nil {
		server.ApiFailed(&w, 201, "更新书签异常"+err.Error())
		return
	}
	util.Logger.Info("update link success" + link.Url)
	server.ApiSuccess(&w, nil)
}
