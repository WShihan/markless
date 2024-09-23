package server

import (
	"markless/util"
	"net/http"
	"strings"
)

func Redirect(w http.ResponseWriter, r *http.Request, route string) {
	if strings.Contains(route, Env.BaseURL) {
		util.Logger.Info("redirect:" + route)
		http.Redirect(w, r, route, http.StatusMovedPermanently)
	} else {
		util.Logger.Info("redirect:" + Env.BaseURL + route)
		http.Redirect(w, r, Env.BaseURL+route, http.StatusMovedPermanently)
	}
}
