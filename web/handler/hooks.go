package handler

import (
	"errors"
	"fmt"
	"markless/model"
	"markless/store"
	"markless/util"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Hooks func(http.Handler) http.Handler

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			util.Logger.Info(r.URL.Path)
			if err := recover(); err != nil {
				err := errors.New(fmt.Sprint(err))
				util.Logger.Fatal(err)
				// 接口返回标准数据
				if strings.Contains(r.URL.Path, "api") {
					ApiFailed(&w, 1, err.Error())
					return
					// 页面接口返回原始页面
				} else {
					SetMsg(&w, err.Error())
					Redirect(w, r, r.Referer())
				}
			}
		}()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		start := time.Now()
		util.Logger.Info(fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
		util.Logger.Info(fmt.Printf("Completed %s in %v", r.URL.Path, time.Since(start)))
	})
}

func Protect(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// 自定义token校验
		xToken := r.Header.Get("X-Token")
		if xToken != "" {
			user := model.User{}
			store.DB.Find(&user, "token = ?", xToken)
			if user.Username != "" {
				util.Logger.Info("Authorization by user token success:" + user.Username)
				r.Header.Set("uid", user.Uid)
				next(w, r, ps)
				return
			} else {
				util.Logger.Info("Authorization by user token failed: invalied token")
			}
		}
		// jwt验证
		authHeader := r.Header.Get("Authorization")
		var jwt = ""
		if authHeader == "" {
			tokenCookie, trr := r.Cookie("markless-token")
			if trr != nil && authHeader == "" {
				util.Logger.Info("Authorization by jwt failed: no token")
				Redirect(w, r, "/login")
				return
			}
			jwt = tokenCookie.Value

		} else {
			jwt = authHeader[len("Bearer "):]
		}

		if jwt == "" {
			util.Logger.Info("Authorization by jwt failed: no token")
			Redirect(w, r, "/login")
			return
		}
		uid, err := util.DecryptAndVerifyJWT(jwt, []byte(Env.HmacSecret), []byte(Env.SecretKey))
		if err != nil {
			util.Logger.Info("Authorization  by jwt failed: validate fails")
			Redirect(w, r, "/login")
			return
		}
		user := model.User{}
		err = store.DB.Find(&user, "uid = ?", uid).Error
		if err != nil {
			util.Logger.Info("Authorization  by jwt failed: validate fails")
			Redirect(w, r, "/login")
			return
		}
		r.Header.Set("uid", uid)
		util.Logger.Info("Authorization by jwt token success:" + jwt)
		next(w, r, ps)
	}
}

func ApplyHooks(h http.Handler, hooks ...Hooks) http.Handler {
	for _, m := range hooks {
		h = m(h)
	}
	return h
}
