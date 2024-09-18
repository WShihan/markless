package api

import (
	"markee/logging"
	"markee/model"
	"markee/store"
	"markee/util"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func TagAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagVal := r.FormValue("tag")
	if tagVal == "" {
		util.Redirect(w, r, "/tags")
		return
	}
	tagArr := strings.Split(tagVal, "&")
	tags := []model.Tag{}

	for _, v := range tagArr {
		tag := model.Tag{Name: strings.Trim(v, " "), UserID: user.ID, CreateTime: time.Now()}
		tags = append(tags, tag)

	}
	store.DB.Model(&user).Association("Tags").Append(&tags)
	err := store.DB.Save(&user).Error
	if err != nil {
		logging.Logger.Error("添加失败" + err.Error())
		model.ApiSuccess(&w, &model.ApiResponse{Msg: err.Error()})
		return
	}
	util.Redirect(w, r, "/tags")
}

func TagDel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	name := params.ByName("name")
	tag := model.Tag{}
	store.DB.Where("name = ? AND user_id = ?", name, user.ID).Find(&tag)
	err := store.DB.Unscoped().Delete(&tag).Error
	if err != nil {
		logging.Logger.Error("删除失败" + err.Error())
		model.ApiSuccess(&w, &model.ApiResponse{Msg: err.Error()})
		return
	}
	util.Redirect(w, r, "/tags")
}
