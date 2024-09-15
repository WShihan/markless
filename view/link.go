package view

import (
	"log/slog"
	"marky/model"
	"marky/store"
	"marky/util"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func LinkCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := r.FormValue("url")
	desc := r.FormValue("desc")
	// category := r.FormValue("category")
	tagName := r.FormValue("tag")
	title, favicon := util.Parse_Webpage(url)
	link := model.Link{}
	user := model.User{}
	tag := model.Tag{}

	store.DB.First(&user)
	store.DB.Find(&link, "url = ?", url)
	store.DB.Find(&tag, "name = ?", tagName)
	if link.Url == "" {
		link.Title = title
		link.Desc = favicon
		link.Desc = desc
		link.CreatedAt = time.Now()
		link.Url = url
		link.UserID = user.ID
		store.DB.Model(&link).Association("Tags").Append(&tag)
		store.DB.Create(&link)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	} else {
		http.Redirect(w, r, "/link/add", http.StatusMovedPermanently)

	}
	slog.Info(title + favicon)

}

func LinkDel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	link := model.Link{}
	store.DB.First(&link, id)
	store.DB.Delete(&link)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
