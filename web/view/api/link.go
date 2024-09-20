package api

import (
	"encoding/json"
	"markless/model"
	"markless/service"
	"markless/store"
	"markless/util"
	"markless/web/handler"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type LinkAddPost struct {
	Url  string `json:"url"`
	Desc string `json:"desc"`
	Tags string `json:"tags"`
	Read bool   `json:"read"`
}

func LinkAddApi(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	postBody := LinkAddPost{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postBody)
	if err != nil {
		panic(err)
	}
	if postBody.Url == "" {
		panic("链接不能为空")
	}
	link, err := service.LinkCreate(user, postBody.Url, postBody.Desc)
	if err != nil {
		panic(err)
	}
	link.Read = postBody.Read
	err = store.DB.Create(&link).Error
	if err != nil {
		panic(err)
	}
	service.LinkAttachTag(user, &link, strings.Split(postBody.Tags, "&"))

	util.Logger.Info("添加书签成功：" + postBody.Url)
	handler.ApiSuccess(&w, &handler.ApiResponse{Msg: "ok", Data: link})

}

func LinkAllApi(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	links := model.Link{}
	store.DB.Where("user_id = ?", user.ID).Find(&links)
	handler.ApiSuccess(&w, &handler.ApiResponse{Msg: "ok", Data: links})
}

func LinkRead(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = true
	err := store.DB.Save(&link).Error
	if err != nil {
		// panic("更新失败" + err.Error())
		handler.ApiSuccess(&w, &handler.ApiResponse{Msg: err.Error()})
		return
	}
	handler.ApiSuccess(&w, &handler.ApiResponse{Msg: "ok", Data: link})
}

func LinkUnread(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = false
	err := store.DB.Save(&link).Error
	if err != nil {
		util.Logger.Error("更新失败" + err.Error())
		handler.ApiSuccess(&w, &handler.ApiResponse{Msg: err.Error()})
		return
	}
	handler.ApiSuccess(&w, &handler.ApiResponse{Msg: "ok", Data: link})
}

func LinkDel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	store.DB.Where("user_id = ?", user.ID).First(&link, id)
	err := store.DB.Unscoped().Delete(&link).Error
	if err != nil {
		util.Logger.Error("删除失败" + err.Error())
		handler.ApiSuccess(&w, &handler.ApiResponse{Msg: err.Error()})
		return
	}
	util.Logger.Info("删除书签成功：" + link.Url)
	handler.ApiSuccess(&w, &handler.ApiResponse{Msg: "ok", Data: link})
}
