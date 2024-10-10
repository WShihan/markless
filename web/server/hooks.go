package server

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
				if strings.Contains(r.URL.Path, "api") {
					// 接口返回标准数据
					switch e := err.(type) {
					case *APIError:
						util.Logger.Error(err)
						ApiFailed(&w, 201, e.Error())
					default:
						util.Logger.Fatal(err)
						ApiFailed(&w, 500, e.Error())
					}
				} else {
					// 页面接口返回原始页面
					SetMsg(&w, err.Error())
					Redirect(w, r, r.Referer())
				}
			}
		}()
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
				http.Redirect(w, r, Env.BaseURL+"/#/login", http.StatusForbidden)
				return
			}
			jwt = tokenCookie.Value

		} else {
			jwt = authHeader[len("Bearer "):]
		}

		if jwt == "" {
			util.Logger.Info("Authorization by jwt failed: no token")
			http.Redirect(w, r, Env.BaseURL+"/#/login", http.StatusForbidden)
			return
		}
		uid, err := util.DecryptAndVerifyJWT(jwt, []byte(Env.HmacSecret), []byte(Env.SecretKey))
		if err != nil {
			util.Logger.Info("Authorization  by jwt failed: validate fails")
			http.Redirect(w, r, Env.BaseURL+"/#/login", http.StatusForbidden)
			return
		}
		user := model.User{}
		err = store.DB.Find(&user, "uid = ?", uid).Error
		if err != nil {
			util.Logger.Info("Authorization  by jwt failed: validate fails")
			http.Redirect(w, r, Env.BaseURL+"/#/login", http.StatusForbidden)
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
