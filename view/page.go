package view

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

	"html/template"
	"log/slog"
	"marky/assets"
	"marky/model"
	"marky/store"
	"path/filepath"

	"marky/util"

	"github.com/julienschmidt/httprouter"
)

type Inject struct {
	Title string
	Data  interface{}
}

func AssetsFinder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	assettype := params.ByName("assettype")
	fileName := params.ByName("filename")
	assetDir := fmt.Sprintf("static/%s/%s", assettype, fileName)
	var fs embed.FS
	if strings.HasSuffix(fileName, ".js") || strings.HasSuffix(fileName, ".map") {
		fs = assets.JS
		w.Header().Set("Content-Type", "text/js")
	} else if strings.HasSuffix(fileName, ".ico") {
		fs = assets.IMG
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
func GetBaseTemplate() *template.Template {
	return template.New("templates/temp.html")
}

func IndexPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/index.html")
	links := []model.Link{}
	store.DB.Find(&links)
	// store.DB.Model(&links).Association("Tags").Find(&[]model.Tag{})
	// store.DB.Model(&links).Association("Tags").Find(&[]model.Tag{})
	inject := Inject{
		Title: "Marky",
		Data:  links,
	}
	tt.ExecuteTemplate(w, "template", inject)
}
func LinkAddPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/link_add.html")
	inject := Inject{
		Title: "添加书签",
	}
	tt.ExecuteTemplate(w, "template", inject)
}
func Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tt, _ := GetBaseTemplate().ParseFS(assets.HTML, "html/template.html", "html/login.html")
	inject := Inject{
		Title: "登录",
	}
	tt.ExecuteTemplate(w, "template", inject)
}

func AdminPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	inject := Inject{
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

func Upload(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 限制上传文件大小
	r.ParseMultipartForm(10 << 20) // 10 MB

	// 从表单中获取文件
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "File uploaded successfully")
}

func ViewPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	picPath := filepath.Join(util.ExcutePath(), "marky_pic", name)
	http.ServeFile(w, r, picPath)
}
