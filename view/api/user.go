package api

import (
	"encoding/json"
	"markee/model"
	"markee/store"
	"markee/util"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func UserAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := model.User{}
	store.DB.Find(&user, "username = ?", username)
	if user.ID == 0 {
		user.Username = username
		user.Password = password
		store.DB.Create(&user)

		res := &model.ApiResponse{Msg: "ok", Data: []interface{}{username}}
		model.ApiSuccess(&w, res)
	} else {
		model.ApiFailed(&w, 1, "用户名已存在")
	}

}

func UserLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := model.User{}
	store.DB.Find(&user, "username = ? AND password = ?", username, password)

	if user.Username != "" {
		token, err := util.CreateJWT(user.Username)
		if err != nil {
			model.ApiFailed(&w, 1, err.Error())
			return
		}

		user.Token = &token
		store.DB.Save(&user)
		cookie := http.Cookie{
			Name:  "markee-token",
			Value: *user.Token,
			Path:  "/",
			// 其他可选字段
			HttpOnly: false, // 使 Cookie 仅通过 HTTP(S) 访问
			Secure:   false, // 在 HTTPS 下设置为 true
			MaxAge:   36000, // Cookie 的有效期（秒）
		}
		// 设置 Cookie
		http.SetCookie(w, &cookie)
		util.Redirect(w, r, "/")
	} else {
		model.ApiFailed(&w, 1, "用户名或密码错误")
	}
}

func UserConfigGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("name")
	var user model.User
	store.DB.Find(&user, "username = ? ", username)
	store.DB.Model(&user).Association("Config").Find(&user)
	if user.Username != "" {
		res := &model.ApiResponse{Msg: "ok", Data: user}
		model.ApiSuccess(&w, res)
	} else {
		model.ApiFailed(&w, 1, "用户名或密码错误")
	}
}

type configForm struct {
	Zoom       int     `form:"zoom"`
	MinZom     int     `form:"minzoom"`
	MaxZoom    int     `form:"maxzoom"`
	Tolorance  float32 `form:"tolerance"`
	Lon        float32 `form:"lon"`
	Lat        float32 `form:"lat"`
	IconSize   int     `form:"iconsize"`
	AutoCenter bool    `form:"autoCenter"`
}

func UserConfigGetUpate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uid := r.Header.Get("uid")
	var user model.User
	store.DB.Find(&user, "uid = ?", uid)
	store.DB.Model(&user).Association("Config").Find(&user)
	if user.Username != "" {
		decoder := json.NewDecoder(r.Body)
		var form configForm
		err := decoder.Decode(&form)
		if err != nil {
			model.ApiFailed(&w, 1, err.Error())
			return
		}

		res := &model.ApiResponse{Msg: "ok", Data: user}
		model.ApiSuccess(&w, res)
	} else {
		model.ApiFailed(&w, 1, "用户名或密码错误")
	}
}
