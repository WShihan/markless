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
		panic(server.APIError{Msg: local.Translate("tip.params.wrong", r.FormValue("lang")), Code: 201})

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
			panic(server.APIError{Msg: local.Translate("tip.password.wrong", user.Lang), Code: 201})
		}
		expires := time.Now().Add(time.Duration(Env.JWTExpire) * time.Minute)
		claims := jwt.MapClaims{
			"uid": user.Uid,
			"exp": expires.Unix(),
		}
		token, err := util.CreateAndEncryptJWT(claims, []byte(Env.HmacSecret), []byte(Env.SecretKey))
		if err != nil {
			panic(server.APIError{Msg: local.Translate("tip.password.wrong", user.Lang), Code: 201})

		}
		data := UserLoginRes{
			Username:    user.Username,
			AccessToken: token,
			Lang:        user.Lang,
			Theme:       user.Theme,
		}
		server.ApiSuccess(&w, &data)
	} else {
		panic(server.APIError{Msg: local.Translate("tip.user.unexist", user.Lang), Code: 201})

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
		panic(server.APIError{Msg: local.Translate("tip.params.wrong", lang), Code: 201})

	}

	if len(post.Username) < 3 || len(post.Password) < 6 {
		msg := server.SetMsg(&w, local.Translate("tip.password.length", r.FormValue("lang")))
		panic(server.APIError{Msg: msg, Code: 201})
	}
	user := model.User{}
	store.DB.Find(&user, "username = ?", post.Username)
	if user.Username == post.Password {
		msg := server.SetMsg(&w, local.Translate("tip.user.existed", r.FormValue("lang")))
		server.ApiFailed(&w, 200, msg)
		return
	} else {
		pass, err := tool.HashMessage(post.Password)
		if err != nil {
			msg := server.SetMsg(&w, local.Translate("msg.failed", r.FormValue("lang")))
			panic(server.APIError{Msg: msg, Code: 201})

		}
		user.Username = post.Username
		user.Password = pass
		user.Lang = lang
		user.Uid = tool.ShortUID(10)
	}
	err = store.DB.Create(&user).Error
	if err != nil {
		panic(server.APIError{Msg: local.Translate("msg.failed", r.FormValue("lang")), Code: 201})
	}
	server.ApiSuccess(&w, nil)

}

func UserInfoGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	if user.Username == "" {
		panic(server.APIError{Msg: local.Translate("tip.user.unexist", r.FormValue("lang")), Code: 201})
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
		panic(server.APIError{Msg: local.Translate("tip.params.wrong", user.Lang), Code: 201})

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
		msg := server.SetMsg(&w, local.Translate("msg.tip.user.unexist", user.Lang))
		panic(server.APIError{Msg: msg, Code: 201})
	}
	user.Token = nil
	store.DB.Save(&user)
	server.ApiSuccess(&w, nil)

}

func UserTokenRefresh(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := store.GetUserByUID(r.Header.Get("uid"))
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("tip.user.unexist", user.Lang))
		panic(server.APIError{Msg: msg, Code: 201})
	}
	tk, err := util.GenerateRandomKey(64)
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("msg.failed", user.Lang))
		panic(server.APIError{Msg: msg, Code: 201})
	}
	user.Token = &tk
	err = store.DB.Save(&user).Error
	if err != nil {
		panic(server.APIError{Msg: local.Translate("msg.failed", user.Lang), Code: 201})

	}
	server.ApiSuccess(&w, nil)

}

func UserUpdatePassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	post := UserPasswordUpdatePost{}
	err := tool.ConvertJSON2Struct(&post, r)
	if err != nil {
		panic(server.APIError{Msg: local.Translate("tip.params.wrong", user.Lang), Code: 201})

	}

	err = tool.ValidateHash(user.Password, post.PasswordOld)
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("tip.password.wrong", user.Lang))
		panic(server.APIError{Msg: msg, Code: 201})

	}
	passwordUpdated, err := tool.HashMessage(post.Password)
	if err != nil {
		msg := server.SetMsg(&w, local.Translate("msg.failed", user.Lang))
		panic(server.APIError{Msg: msg, Code: 201})
	}
	user.Password = passwordUpdated
	store.DB.Save(&user)
	server.ApiSuccess(&w, nil)

}
