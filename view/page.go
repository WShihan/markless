package view

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"html/template"
	"log/slog"
	"markee/assets"
	"markee/model"
	"markee/store"

	"markee/util"

	"github.com/julienschmidt/httprouter"
)

type Search struct {
	Keyword  string
	PrePage  int
	NextPage int
	TagName  string
}
type Inject struct {
	Title string
	Env   model.BaseInjdection
	Search
	Data       interface{}
	TagStastic map[string]int
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
		slog.Error(err.Error())
		return
	}
	w.Write(content)
}

func IndexPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagName := r.URL.Query().Get("tag")
	keywordVal := r.URL.Query().Get("keyword")
	keyword := fmt.Sprintf("%%%s%%", keywordVal)

	pagenum, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	limit := 20
	offset := pagenum * int64(limit)

	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	links := []model.Link{}
	store.DB.Where("Title LIKE ? OR Desc LIKE ?", keyword, keyword).Limit(limit).Offset(int(offset)).Find(&links)
	for i, v := range links {
		tags := []model.Tag{}
		store.DB.Model(&v).Association("Tags").Find(&tags)
		links[i].Tags = tags
	}
	inject := Inject{
		Title:      "markee",
		Env:        Env,
		Search:     Search{Keyword: keywordVal, PrePage: int(pagenum) - 1, NextPage: int(pagenum) + 1, TagName: tagName},
		Data:       links,
		TagStastic: store.TagStat(),
	}
	tt.ExecuteTemplate(w, "template", inject)
}

type EditOption struct {
	Link model.Link
	Tags []model.Tag
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
