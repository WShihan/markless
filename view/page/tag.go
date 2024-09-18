package page

import (
	"net/http"

	"markee/assets"
	"markee/injection"
	"markee/store"
	"markee/util"

	"github.com/julienschmidt/httprouter"
)

func TagsPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/tags.html")
	inject := injection.TagsPage{
		Title: "标签",
		Env:   Env,
		Data:  store.TagStat(user),
	}
	tt.ExecuteTemplate(w, "template", inject)
}
