package api

import (
	"encoding/json"
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

func LinkAllApi(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	links := model.Link{}
	store.DB.Where("user_id = ?", user.ID).Find(&links)
	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok", Data: links})
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
	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok", Data: link})
}

func LinkUnread(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, _ := store.GetLinkByUser(user, int(id))
	link.Read = false
	err := store.DB.Save(&link).Error
	if err != nil {
		util.Logger.Error("update link failed" + err.Error())
		server.ApiSuccess(&w, &server.ApiResponse{Msg: err.Error()})
		return
	}
	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok", Data: link})
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
