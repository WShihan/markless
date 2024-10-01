package api

import (
	"net/http"
	"time"

	"markless/local"
	"markless/model"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/server"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

type UserLoginPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
	Lang        string `json:"lang"`
	Theme       string `json:"theme"`
}

type UserInfo struct {
	Username string  `json:"username"`
	Token    *string `json:"token"`
	Lang     string  `json:"lang"`
	Admin    bool    `json:"admin"`
	Theme    string  `json:"theme"`
}

type UserPasswordUpdatePost struct {
	Password    string `json:"password"`
	PasswordOld string `json:"password_old"`
}

func UserLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	post := UserLoginPost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		util.Logger.Error(err.Error())
		server.ApiFailed(&w, 200, "用户名或密码错误")
		return
	}
	user := model.User{}
	store.DB.Where("username = ?", post.Username).Find(&user)

	if user.Username != "" {
		if user.Password == post.Password {
			// 主要是兼容一开始密码没有保护的情况，新版本第一次登录时会自动hash密码
			pass, err := tool.HashMessage(post.Password)
			if err != nil {
				util.Logger.Error(err.Error())
			}
			user.Password = pass
			store.DB.Save(&user)
		}
		err := tool.ValidateHash(user.Password, post.Password)
		if err != nil {
			server.ApiFailed(&w, 200, "密码错误")
			return
		}
		expires := time.Now().Add(time.Duration(Env.JWTExpire) * time.Minute)
		claims := jwt.MapClaims{
			"uid": user.Uid,
			"exp": expires.Unix(),
		}
		token, err := util.CreateAndEncryptJWT(claims, []byte(Env.HmacSecret), []byte(Env.SecretKey))
		if err != nil {
			server.ApiFailed(&w, 200, "用户名或密码错误")
			return
		}
		data := UserLoginRes{
			Username:    user.Username,
			AccessToken: token,
			Lang:        user.Lang,
			Theme:       user.Theme,
		}
		server.ApiSuccess(&w, &data)
	} else {
		server.ApiFailed(&w, 200, "用户名错误或不存在")
	}
}

func UserRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var lang string
	langCookie, err := r.Cookie("narkless-lang")
	if err != nil {
		lang = local.GetPreferredLanguage(r)
	} else {
		lang = langCookie.Value
	}
	post := UserLoginPost{}
	err = tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		util.Logger.Error(err.Error())
		server.ApiFailed(&w, 200, "用户名或密码错误")
		return
	}

	if len(post.Username) < 3 || len(post.Password) < 6 {
		msg := server.SetMsg(&w, local.Translate("tip.password.length", r.FormValue("lang")))
		server.ApiFailed(&w, 200, msg)
		return
	}
	user := model.User{}
	store.DB.Find(&user, "username = ?", post.Username)
	if user.Username == post.Password {
		msg := server.SetMsg(&w, local.Translate("tip.user.already-exist", r.FormValue("lang")))
		server.ApiFailed(&w, 200, msg)
		return
	} else {
		pass, err := tool.HashMessage(post.Password)
		if err != nil {
			msg := server.SetMsg(&w, local.Translate("msg.failed", r.FormValue("lang")))
			server.ApiFailed(&w, 200, msg)
			return
		}
		user.Username = post.Username
		user.Password = pass
		user.Lang = lang
		user.Uid = tool.ShortUID(10)
	}
	err = store.DB.Create(&user).Error
	if err != nil {
		server.ApiFailed(&w, 200, local.Translate("msg.failed", r.FormValue("lang")))
		return
	}
	server.ApiSuccess(&w, nil)

}

func UserInfoGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	if user.Username == "" {
		server.ApiFailed(&w, 200, "用户不存在")
	}
	if user.Theme == "" {
		user.Theme = "normal"
		store.DB.Save(&user)
	}
	data := UserInfo{
		Username: user.Username,
		Token:    user.Token,
		Lang:     user.Lang,
		Admin:    user.Admin,
		Theme:    user.Theme,
	}
	server.ApiSuccess(&w, data)
}
func UserInfoUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	post := UserInfo{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		server.ApiFailed(&w, 200, "参数错误")
	}
	user.Lang = post.Lang
	user.Theme = post.Theme
	store.DB.Save(&user)

	server.ApiSuccess(&w, nil)
}

func UserEnvGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	server.ApiSuccess(&w, Env)
}

func UserTokenDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("msg.tip.user.not-exist", user.Lang))
		server.ApiFailed(&w, 200, msg)
		return
	}
	user.Token = nil
	store.DB.Save(&user)
	server.ApiSuccess(&w, nil)

}

func UserTokenRefresh(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("tip.user.not-exist", user.Lang))
		server.ApiFailed(&w, 200, msg)
		return
	}
	tk, err := util.GenerateRandomKey(64)
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("msg.failed", user.Lang))
		server.ApiFailed(&w, 200, msg)
		return
	}
	user.Token = &tk
	err = store.DB.Save(&user).Error
	if err != nil {
		server.ApiFailed(&w, 200, "操作失败")
	}
	server.ApiSuccess(&w, nil)

}

func UserUpdatePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	post := UserPasswordUpdatePost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		server.ApiFailed(&w, 200, "参数错误")
	}

	err = tool.ValidateHash(user.Password, post.PasswordOld)
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("tip.password.wrong", user.Lang))
		server.ApiFailed(&w, 200, msg)
		return
	}
	passwordUpdated, err := tool.HashMessage(post.Password)
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("msg.failed", user.Lang))
		server.ApiFailed(&w, 200, msg)
		return
	}
	user.Password = passwordUpdated
	store.DB.Save(&user)
	server.ApiSuccess(&w, nil)

}
