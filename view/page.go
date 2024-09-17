package view

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"html/template"
	"markee/assets"
	"markee/logging"
	"markee/model"
	"markee/store"

	"markee/util"

	"github.com/julienschmidt/httprouter"
)

var (
	LIMIT = 20
)

type Search struct {
	Keyword    string
	PrePage    int
	Page       int
	Limit      int
	NextPage   int
	TagName    string
	ReadStatus int // 0：所有 1:已读 2:未读
	Count      int
}
type Inject struct {
	Title string
	Env   model.BaseInjdection
	Search
	Data       interface{}
	TagStastic map[string]int
}
type EditOption struct {
	Link model.Link
	Tags []model.Tag
}

func GetBaseTemplate() *template.Template {
	funcMap := util.GetFuncMap()
	return template.New("html/template.html").Funcs(funcMap)
}

func AssetsFinder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	assettype := params.ByName("assettype")
	fileName := params.ByName("filename")
	assetDir := fmt.Sprintf("static/%s/%s", assettype, fileName)
	var fs embed.FS
	if strings.HasSuffix(fileName, ".js") || strings.HasSuffix(fileName, ".map") {
		fs = assets.JS
		w.Header().Set("Content-Type", "text/js")
	} else if strings.HasSuffix(fileName, ".png") {
		fs = assets.IMG
		w.Header().Set("Content-Type", "image/x-icon")
	} else if strings.HasSuffix(fileName, ".ico") {
		fs = assets.ICO
		w.Header().Set("Content-Type", "image/x-icon")
	} else {
		fs = assets.CSS
		w.Header().Set("Content-Type", "text/css")
	}
	content, err := fs.ReadFile(assetDir)
	if err != nil {
		logging.Logger.Error(err.Error())
		return
	}
	w.Write(content)
}

func filterLinkByTag(opt *Search) []model.Link {
	rawLinks := []model.Link{}
	offset := opt.Page * opt.Limit
	targetTag := model.Tag{}
	store.DB.Where("name = ?", opt.TagName).Find(&targetTag)
	store.DB.Model(&targetTag).Association("Links").Find(&rawLinks)
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
func filterLinkByKeyword(opt *Search) []model.Link {
	links := []model.Link{}
	offset := opt.Page * opt.Limit
	var count int64
	condition := "%" + opt.Keyword + "%"
	var err error
	if opt.ReadStatus == 0 {
		store.DB.Model(&model.Link{}).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Count(&count)
		err = store.DB.Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
	} else if opt.ReadStatus == 1 {
		store.DB.Model(&model.Link{}).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", true).Count(&count)
		err = store.DB.Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", true).Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
	} else {
		store.DB.Model(&model.Link{}).Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", false).Count(&count)
		err = store.DB.Where("Title LIKE ? OR Desc LIKE ?", condition, condition).Where("read = ?", false).Limit(opt.Limit).Offset(int(offset)).Find(&links).Error
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
	tagName := r.URL.Query().Get("tag")
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	searchOpt := Search{
		PrePage:  int(pagenum) - 1,
		Page:     int(pagenum),
		Limit:    LIMIT,
		NextPage: int(pagenum) + 1,
		TagName:  tagName,
		Keyword:  keyword,
	}

	var links []model.Link
	if tagName != "" {
		links = filterLinkByTag(&searchOpt)
	} else {
		links = filterLinkByKeyword(&searchOpt)
	}
	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	inject := Inject{
		Title:      fmt.Sprintf("所有书签（%d）", searchOpt.Count),
		Env:        Env,
		Search:     searchOpt,
		Data:       links,
		TagStastic: store.TagStat(),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkUnreadPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagName := r.URL.Query().Get("tag")
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	limit := 20
	searchOpt := Search{
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
		links = filterLinkByTag(&searchOpt)
	} else {
		links = filterLinkByKeyword(&searchOpt)
	}

	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	inject := Inject{
		Title:      fmt.Sprintf("未读书签(%d)", searchOpt.Count),
		Env:        Env,
		Search:     searchOpt,
		Data:       links,
		TagStastic: store.TagStat(),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkReadPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagName := r.URL.Query().Get("tag")
	pageVal := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	pagenum, _ := strconv.ParseInt(pageVal, 10, 64)
	limit := 20
	searchOpt := Search{
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
		links = filterLinkByTag(&searchOpt)
	} else {
		links = filterLinkByKeyword(&searchOpt)
	}

	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	inject := Inject{
		Title:      fmt.Sprintf("已读书签(%d)", searchOpt.Count),
		Env:        Env,
		Search:     searchOpt,
		Data:       links,
		TagStastic: store.TagStat(),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkEditPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	link := model.Link{}
	store.DB.Find(&link, id)
	tags := []model.Tag{}
	alltTags := []model.Tag{}
	store.DB.Model(&link).Association("Tags").Find(&tags)
	link.Tags = tags

	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_edit.html")
	user := model.User{}
	store.DB.First(&user)
	store.DB.Model(&user).Association("Tags").Find(&alltTags)
	inject := Inject{
		Title: "编辑书签",
		Env:   Env,
		Data:  EditOption{link, alltTags},
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func LinkAddPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_add.html")
	user := model.User{}
	store.DB.First(&user)
	tags := []model.Tag{}
	store.DB.Model(&user).Association("Tags").Find(&tags)
	inject := Inject{
		Title: "添加书签",
		Env:   Env,
		Data:  tags,
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func TagsPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/tags.html")
	inject := Inject{
		Title: "标签",
		Env:   Env,
		Data:  store.TagStat(),
	}
	tt.ExecuteTemplate(w, "template", inject)
}
func Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 读取嵌入的模板文件
	t, err := template.ParseFS(assets.HTML, "html/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	inject := Inject{
		Env:   Env,
		Title: "登录",
	}
	if err := t.Execute(w, inject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	inject := Inject{
		Env:   Env,
		Title: "首页",
	}
	// 读取嵌入的模板文件
	t, err := template.ParseFS(assets.HTML, "html/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 渲染模板并返回给客户端
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, inject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
