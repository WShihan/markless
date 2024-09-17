package util

import (
	"markee/logging"
	"net/http"
	"strings"
)

func Redirect(w http.ResponseWriter, r *http.Request, route string) {
	if strings.Contains(route, Env.BaseURL) {
		logging.Logger.Info("redirest:" + route)
		http.Redirect(w, r, route, http.StatusMovedPermanently)
	} else {
		logging.Logger.Info("redirest:" + Env.BaseURL + route)
		http.Redirect(w, r, Env.BaseURL+route, http.StatusMovedPermanently)
	}
}
