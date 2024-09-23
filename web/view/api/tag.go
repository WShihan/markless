package api

import (
	"encoding/json"
	"markless/model"
	"markless/store"
	"markless/util"
	"markless/web/server"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func TagDelApi(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagNames := strings.Split(params.ByName("name"), "&")
	for _, v := range tagNames {
		if v == "" {
			continue
		}
		tag := model.Tag{}
		store.DB.Where("name = ? AND user_id = ?", v, user.ID).Find(&tag)
		err := store.DB.Unscoped().Delete(&tag).Error
		if err != nil {
			util.Logger.Error("delete tag failed" + err.Error())
			panic(err)
		}
		util.Logger.Error("deleted tag success" + v)
	}

	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok"})
}

type TagUpdateTitlePost struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TagUpdateName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))

	var data TagUpdateTitlePost
	err := json.NewDecoder(r.Body).Decode(&data) // 读取post数据
	if err != nil {
		panic(err)
	}
	tag := model.Tag{}
	store.DB.Where("id = ? AND user_id = ?", data.ID, user.ID).Find(&tag)
	tag.Name = data.Name
	err = store.DB.Save(&tag).Error
	if err != nil {
		panic(err)
	}

	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok"})
}

type TagUpdateLinkPost struct {
	ID    int   `json:"id"`
	Links []int `json:"links"`
}

func TagUpdateLink(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))

	var data TagUpdateLinkPost
	err := json.NewDecoder(r.Body).Decode(&data) // 读取post数据
	if err != nil {
		panic(err)
	}
	tag := model.Tag{}
	store.DB.Where("id = ? AND user_id = ?", data.ID, user.ID).Find(&tag)
	store.DB.Model(&tag).Association("Links").Clear()

	for _, v := range data.Links {
		link := model.Link{}
		store.DB.Where("id = ? AND user_id = ?", v, user.ID).Find(&link)
		store.DB.Model(&tag).Association("Links").Append(&link)
	}
	err = store.DB.Save(&tag).Error
	if err != nil {
		panic(err)
	}

	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok"})
}
