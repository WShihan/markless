package page

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"markless/injection"
	"markless/local"
	"markless/model"
	"markless/service"
	"markless/store"
	"markless/util"
	"markless/web/assets"
	"markless/web/handler"

	"github.com/julienschmidt/httprouter"
)

var (
	LIMIT = 20
)

func detectPages(count int, limit int) int {
	res := float64(count) / float64(limit)
	return int(math.Ceil(res))
}

func filterLinkByTag(opt *injection.Search, user model.User) []model.Link {
	rawLinks := []model.Link{}
	offset := opt.Page * opt.Limit
	targetTag := model.Tag{}
	store.DB.Where("name = ? and user_id", opt.TagName, user.ID).Find(&targetTag)
	store.DB.Model(&targetTag).Where("user_id = ?", user.ID).Association("Links").Find(&rawLinks)
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
	opt.Pages = detectPages(len(links), opt.Limit)
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
		util.Logger.Error(err.Error())
	}
	opt.Count = int(count)
	opt.Pages = detectPages(opt.Count, opt.Limit)

	// 绑定标签
	for i, v := range links {
		tags := []model.Tag{}
		store.DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}

	return links
}

func IndexPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handler.Redirect(w, r, "/all")
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
		Page: injection.PageInjection{
			Lang:   user.Lang,
			Title:  fmt.Sprintf("%s（%d）", local.Translate("page.all", user.Lang), searchOpt.Count),
			Active: "all",
		},
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
		Page: injection.PageInjection{
			Lang:   user.Lang,
			Title:  fmt.Sprintf("%s（%d）", local.Translate("page.unread", user.Lang), searchOpt.Count),
			Active: "unread",
		},
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
		Page: injection.PageInjection{
			Lang:   user.Lang,
			Title:  fmt.Sprintf("%s（%d）", local.Translate("page.read", user.Lang), searchOpt.Count),
			Active: "read",
		},
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
		Page: injection.PageInjection{
			Lang:   user.Lang,
			Title:  local.Translate("page.edit-link", user.Lang),
			Active: "",
		},
		Env:  Env,
		Data: injection.LinkEditInjection{Link: link, Tags: alltTags},
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkAddPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_add.html")
	tags := []model.Tag{}
	store.DB.Model(&user).Where("user_id = ?", user.ID).Association("Tags").Find(&tags)
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Title:  local.Translate("page.link-find", user.Lang),
			Lang:   user.Lang,
			Active: "link-find",
		},
		Env:  Env,
		Data: tags,
	}
	tt.ExecuteTemplate(w, "template", inject)
}
func LinkArchViewPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	linkID := params.ByName("id")

	link := model.Link{}
	store.DB.Preload("Archive").Find(&link, linkID)
	store.DB.Preload("Tags").Find(&link, linkID)

	tt, _ := util.GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_archive.html")
	inject := injection.LinkPage{
		Page: injection.PageInjection{
			Title:  local.Translate("page.link.archive.title", user.Lang),
			Lang:   user.Lang,
			Active: "",
		},
		Env:  Env,
		Data: link,
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	webURL := r.FormValue("url")
	tagNames := strings.Split(r.FormValue("tags"), "&")
	desc := r.FormValue("desc")
	if webURL == "" {
		panic(local.Translate("tip.link.empty", user.Lang))
	}

	link, err := service.LinkCreate(user, webURL, desc)
	if err != nil {
		panic(err)
	}
	service.LinkAttachTag(user, &link, tagNames)
	handler.Redirect(w, r, "/")

}

func LinkUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id, _ := strconv.ParseInt(params.ByName("id"), 10, 64)
	link, err := store.GetLinkByUser(user, int(id))
	if err != nil {
		panic(local.Translate("tip.link.not-exist", user.Lang))
	} else {
		link.Title = r.FormValue("title")
		link.Desc = r.FormValue("desc")
		link.Url = r.FormValue("url")
		link.Icon = r.FormValue("icon")
		tagArr := strings.Split(r.FormValue("tags"), "&")
		tags := []model.Tag{}
		for _, v := range tagArr {
			if v == "" {
				continue
			}
			tag := model.Tag{}
			store.DB.Find(&tag, "name = ?", strings.Trim(v, " "))
			tags = append(tags, tag)
		}
		store.DB.Model(&link).Association("Tags").Append(&tags)
		err = store.DB.Save(&link).Error
		if err != nil {
			panic(local.Translate("msg.failed", user.Lang) + err.Error())
		}
	}
	handler.Redirect(w, r, "/")
}
func LinkUpdateArchive(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	id := params.ByName("id")
	link := model.Link{}
	store.DB.Preload("Archive").Where("user_id = ? AND id = ?", user.ID, id).First(&link)
	if link.Url == "" {
		handler.Redirect(w, r, r.Referer())
		return
	}
	pageInfo, err := util.ParsePage(link.Url)
	if err != nil {
		handler.Redirect(w, r, r.Referer())
		return
	}
	if link.Archive != nil {
		link.Archive.Content = pageInfo.Content
		link.Archive.UpdateTime = time.Now()
	} else {
		link.Archive = &model.Archive{
			LinkID:     link.ID,
			Content:    pageInfo.Content,
			UpdateTime: time.Now(),
		}
	}

	err = store.DB.Save(&link.Archive).Error
	if err != nil {
		handler.Redirect(w, r, r.Referer())
		return
	}
	util.Logger.Info("update link success" + link.Url)
	handler.Redirect(w, r, r.Referer())
}

func LinkMarkAllAsRead(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	link := model.Link{}
	store.DB.Model(&link).Where("user_id = ?", user.ID).Updates(model.Link{Read: true})
	msg := local.Translate("tip.link.mark-all-read", user.Lang)
	util.Logger.Info(msg)
	handler.SetMsg(&w, msg)
	handler.Redirect(w, r, r.Referer())

}

func LinkMarkAllAsUnread(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, _ := store.GetUserByUID(r.Header.Get("uid"))
	links := []model.Link{}
	store.DB.Model(&user).Association("Links").Find(&links)
	for _, v := range links {
		v.Read = false
		store.DB.Save(&v)
	}
	msg := local.Translate("tip.link.mark-all-unread", user.Lang)
	util.Logger.Info(msg)
	handler.SetMsg(&w, msg)
	handler.Redirect(w, r, r.Referer())
}
