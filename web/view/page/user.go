package page

import (
	"net/http"

	"html/template"
	"markless/assets"
	"markless/injection"
	"markless/model"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/handler"

	"github.com/julienschmidt/httprouter"
)

func LoginPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t, err := template.ParseFS(assets.HTML, "html/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	inject := injection.UserLoginPage{
		Env:   Env,
		Title: "登录",
	}
	if err := t.Execute(w, inject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t, err := template.ParseFS(assets.HTML, "html/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	inject := injection.UserLoginPage{
		Env:   Env,
		Title: "注册",
	}
	if err := t.Execute(w, inject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SettingPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/setting.html")
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Active: "setting",
			Title:  "设置",
		},
		Env:  Env,
		Data: user,
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func UserLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := model.User{}
	store.DB.Where("username = ? AND password = ?", username, password).Find(&user)

	if user.Username != "" {
		if user.Password != password {
			handler.SetMsg(&w, "密码错误")
			handler.Redirect(w, r, "/login")
			return
		}
		token, err := util.CreateJWT(user.Uid)
		if err != nil {
			handler.SetMsg(&w, "用户名或密码错误")
			handler.Redirect(w, r, "/login")
			return
		}

		store.DB.Save(&user)
		cookie := http.Cookie{
			Name:  "markless-token",
			Value: token,
			Path:  "/",
			// 其他可选字段
			HttpOnly: false,       // 使 Cookie 仅通过 HTTP(S) 访问
			Secure:   false,       // 在 HTTPS 下设置为 true
			MaxAge:   60 * 60 * 1, // Cookie 的有效期（秒）
		}
		// 设置 Cookie
		http.SetCookie(w, &cookie)
		handler.Redirect(w, r, "/")
	} else {
		handler.SetMsg(&w, "用户不存在")
		handler.Redirect(w, r, "/login")
	}
}

func UserRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password-confirm")

	if password != passwordConfirm {
		handler.SetMsg(&w, "新密码不一致")
		handler.Redirect(w, r, "/register")
		return
	}
	if len(username) < 3 || len(password) < 6 {
		handler.SetMsg(&w, "用户名或密码长度不正确")
		handler.Redirect(w, r, "/register")
		return
	}
	user := model.User{}
	store.DB.Find(&user, "username = ?", username)
	if user.Username == username {
		handler.SetMsg(&w, "用户名已存在")
		handler.Redirect(w, r, "/register")
		return
	} else {
		user.Username = username
		user.Password = password
		user.Uid = tool.ShortUID(10)
	}
	store.DB.Create(&user)
	handler.Redirect(w, r, "/login")

}

func UserChangePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	passwordOld := r.FormValue("password-old")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password-confirm")

	if password != passwordConfirm {
		handler.SetMsg(&w, "新密码不一致")
		handler.Redirect(w, r, "/setting")
		return
	}
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		handler.SetMsg(&w, "用户不存在")
		handler.Redirect(w, r, "/setting")
		return
	}

	if user.Password != passwordOld {
		handler.SetMsg(&w, "原始密码错误")
		handler.Redirect(w, r, "/setting")
		return
	}
	user.Password = password
	store.DB.Save(&user)
	handler.SetMsg(&w, "密码修改成功")
	handler.Redirect(w, r, "/setting")

}

func UserTokenAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		handler.SetMsg(&w, "用户不存在")
		handler.Redirect(w, r, "/setting")
		return
	}
	tk, err := util.GenerateRandomKey(64)
	if err != nil {
		handler.SetMsg(&w, "生成token失败")
		handler.Redirect(w, r, "/setting")
		return
	}
	user.Token = &tk
	handler.SetMsg(&w, "创建成功")

	store.DB.Save(&user)
	handler.Redirect(w, r, "/setting")

}

func UserTokenDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		handler.SetMsg(&w, "用户不存在")
		handler.Redirect(w, r, "/setting")
		return
	}
	user.Token = nil
	store.DB.Save(&user)
	handler.SetMsg(&w, "删除成功")
	handler.Redirect(w, r, "/setting")

}
