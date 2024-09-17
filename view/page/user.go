package page

import (
	"net/http"

	"html/template"
	"markee/assets"
	"markee/injection"

	"github.com/julienschmidt/httprouter"
)

func Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 读取嵌入的模板文件
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
