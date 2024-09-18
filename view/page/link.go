package page

import (
	"fmt"
	"net/http"
	"strconv"

	"markee/assets"
	"markee/injection"
	"markee/logging"
	"markee/model"
	"markee/store"
	"markee/util"

	"github.com/julienschmidt/httprouter"
)

var (
	LIMIT = 20
)

func filterLinkByTag(opt *injection.Search, user model.User) []model.Link {
	rawLinks := []model.Link{}
	offset := opt.Page * opt.Limit
	targetTag := model.Tag{}
	store.DB.Where("name = ? and user_id", opt.TagName, user.ID).Find(&targetTag)
	store.DB.Model(&targetTag).Where("user_id = ?", user.ID).Preload("Tags").Association("Links").Find(&rawLinks)
	links := rawLinks[:0]
	if opt.ReadStatus == 0 {
		links = append(links, rawLinks...)
	} else if opt.ReadStatus == 1 {
		for _, v := range rawLinks {
			if v.Read {
				links = append(links, v)
			}
		}
	} else if opt.ReadStatus == 2 {
		for _, v := range rawLinks {
			if !v.Read {
				links = append(links, v)
			}
		}
	}
	opt.Count = len(links)
	opt.Keyword = "#" + opt.TagName

	if int(offset) < len(links) {
		end := int(offset) + opt.Limit
		if end > len(links) {
			end = len(links)
		}
		links = links[offset:end]
	}

	// 绑定标签
	for i, v := range links {
		tags := []model.Tag{}
		store.DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}

	return links
}

func filterLinkByKeyword(opt *injection.Search, user model.User) []model.Link {
	links := []model.Link{}
	offset := opt.Page * opt.Limit
	var count int64
	condition := "%" + opt.Keyword + "%"
	var err error
	if opt.ReadStatus == 0 {
		store.DB.Model(&model.Link{}).Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Count(&count)
		err = store.DB.Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Order("created_at DESC").Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
	} else if opt.ReadStatus == 1 {
		store.DB.Model(&model.Link{}).Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", true).Count(&count)
		err = store.DB.Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", true).Order("created_at DESC").Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
	} else {
		store.DB.Model(&model.Link{}).Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", false).Count(&count)
		err = store.DB.Where("user_id = ?", user.ID).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", false).Limit(opt.Limit).Order("created_at DESC").Offset(int(offset)).Find(&links).Error
	}
	if err != nil {
		logging.Logger.Error(err.Error())
	}
	opt.Count = int(count)

	// 绑定标签
	for i, v := range links {
		tags := []model.Tag{}
		store.DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}

	return links
}

func IndexPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	util.Redirect(w, r, "/unread")
}

func LinkAllPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagName := r.URL.Query().Get("tag")
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	searchOpt := injection.Search{
		PrePage:  int(pagenum) - 1,
		Page:     int(pagenum),
		Limit:    LIMIT,
		NextPage: int(pagenum) + 1,
		TagName:  tagName,
		Keyword:  keyword,
	}

	var links []model.Link
	if tagName != "" {
		links = filterLinkByTag(&searchOpt, user)
	} else {
		links = filterLinkByKeyword(&searchOpt, user)
	}
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	inject := injection.LinkPage{
		Title:      fmt.Sprintf("所有书签（%d）", searchOpt.Count),
		Env:        Env,
		Search:     searchOpt,
		Data:       links,
		TagStastic: store.TagStat(user),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkUnreadPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagName := r.URL.Query().Get("tag")
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	limit := 20
	searchOpt := injection.Search{
		PrePage:    int(pagenum) - 1,
		Page:       int(pagenum),
		Limit:      limit,
		NextPage:   int(pagenum) + 1,
		TagName:    tagName,
		Keyword:    keyword,
		ReadStatus: 2,
	}

	var links []model.Link
	if tagName != "" {
		links = filterLinkByTag(&searchOpt, user)
	} else {
		links = filterLinkByKeyword(&searchOpt, user)
	}

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	inject := injection.LinkPage{
		Title:      fmt.Sprintf("未读书签(%d)", searchOpt.Count),
		Env:        Env,
		Search:     searchOpt,
		Data:       links,
		TagStastic: store.TagStat(user),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkReadPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tagName := r.URL.Query().Get("tag")
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	limit := 20
	searchOpt := injection.Search{
		PrePage:    int(pagenum) - 1,
		Page:       int(pagenum),
		Limit:      limit,
		NextPage:   int(pagenum) + 1,
		TagName:    tagName,
		Keyword:    keyword,
		ReadStatus: 1,
	}

	var links []model.Link
	if tagName != "" {
		links = filterLinkByTag(&searchOpt, user)
	} else {
		links = filterLinkByKeyword(&searchOpt, user)
	}

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	inject := injection.LinkPage{
		Title:      fmt.Sprintf("已读书签(%d)", searchOpt.Count),
		Env:        Env,
		Search:     searchOpt,
		Data:       links,
		TagStastic: store.TagStat(user),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkEditPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id := params.ByName("id")
	link := model.Link{}
	store.DB.Find(&link, id)
	tags := []model.Tag{}
	alltTags := []model.Tag{}
	store.DB.Model(&link).Where("user_id = ?", user.ID).Association("Tags").Find(&tags)
	link.Tags = tags

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_edit.html")
	store.DB.Model(&user).Where("user_id = ?", user.ID).Association("Tags").Find(&alltTags)
	inject := injection.LinkPage{
		Title: "编辑书签",
		Env:   Env,
		Data:  injection.LinkEditInjection{Link: link, Tags: alltTags},
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkAddPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_add.html")
	tags := []model.Tag{}
	store.DB.Model(&user).Where("user_id = ?", user.ID).Association("Tags").Find(&tags)
	inject := injection.LinkPage{
		Title: "添加书签",
		Env:   Env,
		Data:  tags,
	}
	tt.ExecuteTemplate(w, "template", inject)
}
