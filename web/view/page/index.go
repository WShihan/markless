package page

import (
	"html/template"
	"net/http"

	"markless/web/assets"

	"github.com/julienschmidt/httprouter"
)

func IndexPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t, _ := template.ParseFS(assets.HTML, "index.html")
	t.Execute(w, nil)
}
