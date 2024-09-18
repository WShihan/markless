package hooks

import (
	"fmt"
	"markee/logging"
	"markee/model"
	"markee/store"
	"markee/util"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Hooks func(http.Handler) http.Handler

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // 允许所有域名
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // 允许的请求方法
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // 允许的请求头
		start := time.Now()
		logging.Logger.Info(fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
		logging.Logger.Info(fmt.Printf("Completed %s in %v", r.URL.Path, time.Since(start)))
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
				logging.Logger.Info("Authorization by user token success:" + user.Username)
				r.Header.Set("uid", user.Uid)
				next(w, r, ps)
				return
			} else {
				logging.Logger.Info("Authorization by user token failed: invalied token")
			}
		}
		// jwt验证
		authHeader := r.Header.Get("Authorization")
		var jwt = ""
		if authHeader == "" {
			tokenCookie, trr := r.Cookie("markee-token")
			if trr != nil && authHeader == "" {
				logging.Logger.Info("Authorization by jwt failed: no token")
				util.Redirect(w, r, "/login")
				return
			}
			jwt = tokenCookie.Value

		} else {
			jwt = authHeader[len("Bearer "):]
		}

		if jwt == "" {
			logging.Logger.Info("Authorization by jwt failed: no token")
			util.Redirect(w, r, "/login")
			return
		}
		uid, err := util.ValidateJWT(jwt)
		if err != nil {
			logging.Logger.Info("Authorization  by jwt failed: validate fails")
			util.Redirect(w, r, "/login")
			return
		}
		user := model.User{}
		err = store.DB.Find(&user, "uid = ?", uid).Error
		if err != nil {
			logging.Logger.Info("Authorization  by jwt failed: validate fails")
			util.Redirect(w, r, "/login")
			return
		}
		r.Header.Set("uid", uid)
		logging.Logger.Info("Authorization by jwt token success:" + jwt)
		next(w, r, ps)
	}
}

func ApplyMiddleware(h http.Handler, hooks ...Hooks) http.Handler {
	for _, m := range hooks {
		h = m(h)
	}
	return h
}
