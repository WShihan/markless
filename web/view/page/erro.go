package page

import (
	"html/template"
	"markless/assets"
	"markless/injection"
	"markless/local"
	"markless/store"
	"markless/util"
	"net/http"
)

type ErrorMsg struct {
	Msg  string
	Desc string
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	var tt *template.Template
	if user.Username == "" {
		tt, _ = util.GetBaseTemplate().ParseFS(assets.HTML, "html/template_unlogin.html", "html/error.html")

	} else {
		tt, _ = util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/error.html")

	}
	w.WriteHeader(http.StatusNotFound)
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Title:  local.Translate("tip.404.title", user.Lang),
			Active: "",
		},
		Env:  Env,
		Data: ErrorMsg{Msg: local.Translate("tip.404.title", user.Lang), Desc: local.Translate("tip.404.desctiption", user.Lang)},
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func MthodNotAllowedPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	lang := local.GetPreferredLanguage(r)
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template_unlogin.html", "html/error.html")
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Title:  local.Translate("tip.method-not-allowed.title", lang),
			Active: "",
		},
		Env:  Env,
		Data: ErrorMsg{Msg: local.Translate("tip.method-not-allowed.title", lang), Desc: local.Translate("tip.method-not-allowed.desctiption", lang)},
	}
	tt.ExecuteTemplate(w, "template", inject)
}
