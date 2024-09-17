package page

import (
	"net/http"

	"html/template"
	"markee/assets"
	"markee/injection"
	"markee/model"
	"markee/store"
	"markee/util"

	"github.com/julienschmidt/httprouter"
)

func LoginPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t, err := template.ParseFS(assets.HTML, "html/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	inject := injection.UserLoginPage{
		Env:   Env,
		Title: "登录",
	}
	if err := t.Execute(w, inject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SettingPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/setting.html")
	user := model.User{}
	store.DB.First(&user)
	inject := injection.LinkPage{
		Title: "设置",
		Env:   Env,
		Data:  user,
	}
	tt.ExecuteTemplate(w, "template", inject)
}
