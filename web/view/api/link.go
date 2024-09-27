package api

import (
	"encoding/json"
	"markless/injection"
	"markless/local"
	"markless/model"
	"markless/service"
	"markless/store"
	"markless/util"
	"markless/web/server"
	"net/http"
	"strconv"
	"strings"

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

func LinkAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")
	tagName := ""
	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	pagenum = pagenum - 1

	var links []model.Link
	if strings.Contains(keyword, "#") {
		tagName = strings.Replace(keyword, "#", "", 1)
	}
	searchOpt := injection.Search{
		PrePage:  int(pagenum) - 1,
		Page:     int(pagenum),
		Limit:    LIMIT,
		NextPage: int(pagenum) + 2,
		TagName:  tagName,
		Keyword:  keyword,
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

type LinkExistPost struct {
	Url string `json:"url"`
}

func ConvertJSON2Struct(data interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}
	return nil
}

func LinkExist(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	postData := LinkExistPost{}
	err := ConvertJSON2Struct(&postData, r)
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
