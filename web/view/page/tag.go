package page

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"markless/assets"
	"markless/injection"
	"markless/local"
	"markless/model"
	"markless/store"
	"markless/util"
	"markless/web/handler"

	"github.com/julienschmidt/httprouter"
)

func TagsPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/tags.html")
	stat := store.TagStat(user)
	inject := injection.TagsPage{
		Page: injection.PageInjection{
			Title:  fmt.Sprintf("%s（%d）", local.Translate("page.tags", user.Lang), len(stat)),
			Active: "tags",
		},
		Env:  Env,
		Data: stat,
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func TagAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagVal := r.FormValue("tag")
	if tagVal == "" {
		handler.Redirect(w, r, "/tags")
		return
	}
	tagArr := strings.Split(tagVal, "&")
	tags := []model.Tag{}

	for _, v := range tagArr {
		tag := model.Tag{Name: strings.Trim(v, " "), UserID: user.ID, CreateTime: time.Now()}
		store.DB.Where("name = ? AND user_id = ?", tag.Name, user.ID).Find(&tag)
		if tag.ID != 0 {
			util.Logger.Info(local.Translate("tip.tag.unique", user.Lang) + tag.Name)
			continue
		}
		tags = append(tags, tag)

	}
	store.DB.Model(&user).Association("Tags").Append(&tags)
	err := store.DB.Save(&user).Error
	if err != nil {
		util.Logger.Error(local.Translate("msg.failed", user.Lang) + err.Error())
		return
	}
	handler.Redirect(w, r, "/tags")
}

type TagEditPageData struct {
	Tag       model.Tag
	Applied   []model.Link
	Unapplied []model.Link
}

func TagEditPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id := params.ByName("id")
	tag := model.Tag{}
	err := store.DB.Where("id = ? AND user_id = ?", id, user.ID).Find(&tag).Error
	if err != nil {
		panic(err)
	}
	links := []model.Link{}
	store.DB.Model(&tag).Where("user_id = ?", user.ID).Association("Links").Find(&links)

	allLinks := []model.Link{}
	store.DB.Model(&user).Where("user_id = ?", user.ID).Association("Links").Find(&allLinks)
	unapplied := []model.Link{}

	for _, v := range allLinks {
		applied := false
		for _, vv := range links {
			if v.ID == vv.ID {
				applied = true
				break
			}
		}
		if applied {
			continue
		} else {
			unapplied = append(unapplied, v)
		}
	}

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/tag_edit.html")
	inject := injection.TagsPage{
		Page: injection.PageInjection{
			Title:  local.Translate("page.edit-tag", user.Lang),
			Active: "",
		},
		Env:  Env,
		Data: TagEditPageData{Tag: tag, Applied: links, Unapplied: unapplied},
	}
	tt.ExecuteTemplate(w, "template", inject)
}
