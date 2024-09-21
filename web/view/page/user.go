package page

import (
	"net/http"

	"markless/assets"
	"markless/injection"
	"markless/local"
	"markless/model"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/handler"

	"github.com/julienschmidt/httprouter"
)

func LoginPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	lang := local.GetPreferredLanguage(r)

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template_unlogin.html", "html/login.html")
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Title:  local.Translate("page.login", lang),
			Active: "",
			Lang:   lang,
		},
		Env: Env,
	}
	tt.ExecuteTemplate(w, "template", inject)

}

func RegisterPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	lang := local.GetPreferredLanguage(r)

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template_unlogin.html", "html/register.html")
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Lang:   lang,
			Title:  local.Translate("page.register", lang),
			Active: "",
		},
		Env: Env,
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func SettingPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/setting.html")
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Active: "setting",
			Lang:   user.Lang,
			Title:  local.Translate("page.setting", user.Lang),
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
			handler.SetMsg(&w, local.Translate("page.login.error.password", user.Lang))
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
		msg := local.Translate("tip.user.not-exist", user.Lang)
		handler.SetMsg(&w, msg)
		handler.Redirect(w, r, "/login")
	}
}

func UserRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password-confirm")

	if password != passwordConfirm {
		handler.SetMsg(&w, local.Translate("tip.password.not-match", r.FormValue("lang")))
		handler.Redirect(w, r, "/register")
		return
	}
	if len(username) < 3 || len(password) < 6 {
		handler.SetMsg(&w, local.Translate("tip.password.length", r.FormValue("lang")))
		handler.Redirect(w, r, "/register")
		return
	}
	user := model.User{}
	store.DB.Find(&user, "username = ?", username)
	if user.Username == username {
		handler.SetMsg(&w, local.Translate("tip.user.already-exist", r.FormValue("lang")))
		handler.Redirect(w, r, "/register")
		return
	} else {
		user.Username = username
		user.Password = password
		user.Lang = local.GetPreferredLanguage(r)
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
		handler.SetMsg(&w, local.Translate("tip.password.not-match", r.FormValue("lang")))
		handler.Redirect(w, r, "/setting")
		return
	}
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		handler.SetMsg(&w, local.Translate("tip.user.not-exist", user.Lang))
		handler.Redirect(w, r, "/setting")
		return
	}

	if user.Password != passwordOld {
		handler.SetMsg(&w, local.Translate("tip.password.wrong", user.Lang))
		handler.Redirect(w, r, "/setting")
		return
	}
	user.Password = password
	store.DB.Save(&user)
	handler.SetMsg(&w, local.Translate("msg.updated", user.Lang))
	handler.Redirect(w, r, "/setting")

}

func UserTokenAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		handler.SetMsg(&w, local.Translate("tip.user.not-exist", user.Lang))
		handler.Redirect(w, r, "/setting")
		return
	}
	tk, err := util.GenerateRandomKey(64)
	if err != nil {
		handler.SetMsg(&w, local.Translate("msg.failed", user.Lang))
		handler.Redirect(w, r, "/setting")
		return
	}
	user.Token = &tk
	handler.SetMsg(&w, local.Translate("msg.created", user.Lang))

	store.DB.Save(&user)
	handler.Redirect(w, r, "/setting")

}

func UserTokenDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		handler.SetMsg(&w, local.Translate("msg.tip.user.not-exist", user.Lang))
		handler.Redirect(w, r, "/setting")
		return
	}
	user.Token = nil
	store.DB.Save(&user)
	handler.SetMsg(&w, local.Translate("msg.deleted", user.Lang))
	handler.Redirect(w, r, "/setting")

}

func UserBasicUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	lang := r.FormValue("lang")
	user.Lang = lang
	store.DB.Save(&user)
	handler.SetMsg(&w, local.Translate("msg.success", user.Lang))
	handler.Redirect(w, r, "/setting")

}
