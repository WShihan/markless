package api

import (
	"encoding/json"
	"markless/local"
	"markless/model"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/server"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type TagUpdateTitlePost struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TagUpdateLinkPost struct {
	ID    int   `json:"id"`
	Links []int `json:"links"`
}

type TagAttachLinksPost struct {
	Links []int  `json:"links"`
	Tag   string `json:"tag"`
}

type TagAddPost struct {
	Names []string `json:"names"`
}

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
			panic(server.APIError{
				Msg:  local.Translate("tip.params.wrong", user.Lang) + err.Error(),
				Code: 201},
			)
		}
		util.Logger.Error(local.Translate("tip.tag.deleted", user.Lang) + v)
	}

	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok"})
}

func TagAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	err := store.DB.Preload("Tags").Where("id = ?", user.ID).Find(&user).Error
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})

	}
	server.ApiSuccess(&w, &user.Tags)
}

func TagUpdateName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))

	var data TagUpdateTitlePost
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})

	}
	tag := model.Tag{}
	store.DB.Where("id = ? AND user_id = ?", data.ID, user.ID).Find(&tag)
	tag.Name = data.Name
	err = store.DB.Save(&tag).Error
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}

	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok"})
}

func TagAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	var data TagAddPost
	err := tool.ConvertJSON2Struct(&data, r)
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}

	for _, v := range data.Names {
		if v == "" {
			continue
		}
		tag := model.Tag{}
		store.DB.Where("name = ? AND user_id = ?", v, user.ID).Find(&tag)
		if tag.Name == "" {
			tag = model.Tag{Name: v, UserID: user.ID, CreateTime: time.Now()}
			store.DB.Create(&tag)
			store.DB.Model(&user).Association("Tags").Append(&tag)
		}
	}
	err = store.DB.Save(&user).Error
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}
	server.ApiSuccess(&w, nil)

}

func TagUpdateLink(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))

	var data TagUpdateLinkPost
	err := json.NewDecoder(r.Body).Decode(&data) // 读取post数据
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
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
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}

	server.ApiSuccess(&w, &server.ApiResponse{Msg: "ok"})
}

func TagRelatedLinks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagVal := params.ByName("name")
	tag := model.Tag{}
	err := store.DB.Where("name = ? AND user_id = ?", tagVal, user.ID).Find(&tag).Error
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}
	store.DB.Preload("Links").Find(&tag)
	server.ApiSuccess(&w, &tag.Links)

}

func AttachLinks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	post := TagAttachLinksPost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}
	tag := model.Tag{}
	err = store.DB.Where("name = ? AND user_id = ?", post.Tag, user.ID).Find(&tag).Error
	if err != nil {
		panic(server.APIError{Msg: err.Error(), Code: 201})
	}
	store.DB.Model(&tag).Association("Links").Clear()
	for _, v := range post.Links {
		link := model.Link{}
		err := store.DB.Where("id = ? AND user_id = ?", v, user.ID).Find(&link).Error
		if err != nil {
			panic(server.APIError{Msg: err.Error(), Code: 201})
		}
		store.DB.Model(&tag).Association("Links").Append(&link)
	}
	store.DB.Save(&tag)
	server.ApiSuccess(&w, nil)

}

func TagStastic(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	info := store.TagStat(user)
	server.ApiSuccess(&w, info)
}
