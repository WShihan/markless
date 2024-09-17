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
				next(w, r, ps)
				return
			}
		}
		authHeader := r.Header.Get("Authorization")
		var jwt = ""
		if authHeader == "" {
			tokenCookie, trr := r.Cookie("markee-token")
			if trr != nil && authHeader == "" {
				util.Redirect(w, r, "/login")
				return
			}
			jwt = tokenCookie.Value

		} else {
			jwt = authHeader[len("Bearer "):]
		}

		if jwt == "" {
			util.Redirect(w, r, "/login")
			return
		}

		_, err := util.ValidateJWT(jwt)
		if err != nil {
			util.Redirect(w, r, "/login")
			return
		}
		next(w, r, ps)
	}
}

func ApplyMiddleware(h http.Handler, hooks ...Hooks) http.Handler {
	for _, m := range hooks {
		h = m(h)
	}
	return h
}
