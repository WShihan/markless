package api

import (
	"markee/model"
	"markee/store"
	"markee/tool"
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
			tool.SetMsg(&w, "用户名或密码错误")
			util.Redirect(w, r, "/login")
			return
		}

		store.DB.Save(&user)
		cookie := http.Cookie{
			Name:  "markee-token",
			Value: token,
			Path:  "/",
			// 其他可选字段
			HttpOnly: false,       // 使 Cookie 仅通过 HTTP(S) 访问
			Secure:   false,       // 在 HTTPS 下设置为 true
			MaxAge:   60 * 60 * 1, // Cookie 的有效期（秒）
		}
		// 设置 Cookie
		http.SetCookie(w, &cookie)
		util.Redirect(w, r, "/")
	} else {
		tool.SetMsg(&w, "用户不存在")
		util.Redirect(w, r, "/login")
	}
}

func UserChangePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	passwordOld := r.FormValue("password-old")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password-confirm")

	if password != passwordConfirm {
		tool.SetMsg(&w, "新密码不一致")
		util.Redirect(w, r, "/setting")
		return
	}
	user := model.User{}
	store.DB.Find(&user, "password = ?", passwordOld)
	if user.ID == 0 {
		tool.SetMsg(&w, "原始密码错误")
		util.Redirect(w, r, "/setting")
		return
	}
	user.Password = password
	store.DB.Save(&user)
	tool.SetMsg(&w, "密码修改成功")
	util.Redirect(w, r, "/setting")

}

func UserTokenAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := model.User{}
	store.DB.First(&user)
	if user.ID != 0 {
		tk, err := util.GenerateRandomKey(64)
		if err != nil {
			tool.SetMsg(&w, "生成token失败")
			util.Redirect(w, r, "/setting")
			return
		}
		user.Token = &tk
		tool.SetMsg(&w, "创建成功")

		store.DB.Save(&user)
	} else {
		tool.SetMsg(&w, "用户不存在")
	}
	util.Redirect(w, r, "/setting")

}

func UserTokenDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := model.User{}
	store.DB.First(&user)
	user.Token = nil
	store.DB.Save(&user)
	tool.SetMsg(&w, "删除成功")
	util.Redirect(w, r, "/setting")

}
