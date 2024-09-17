package view

import (
	"markee/logging"
	"markee/model"
	"net/http"
	"strings"
)

var (
	BASE_URL = "/markee"
	Env      model.BaseInjdection
)

func Redirect(w http.ResponseWriter, r *http.Request, route string) {
	if strings.Contains(route, Env.BaseURL) {
		logging.Logger.Info("redirest:" + route)
		http.Redirect(w, r, route, http.StatusMovedPermanently)
	} else {
		logging.Logger.Info("redirest:" + Env.BaseURL + route)
		http.Redirect(w, r, Env.BaseURL+route, http.StatusMovedPermanently)
	}
}

func LoadApi(router *model.RouterWithPrefix) {
	router.POST("/user/login", UserLogin)
	router.POST("/link/add", Protect(LinkAdd))
	router.POST("/link/update/:id", Protect(LinkUpdate))
	router.POST("/tag/add", Protect(TagAdd))

	router.GET("/link/read/:id", Protect(LinkRead))
	router.GET("/link/unread/:id", Protect(LinkUnread))

}

func LoadPage(router *model.RouterWithPrefix, env model.BaseInjdection) {
	Env = env
	router.GET("/login", Login)
	router.GET("/", Protect(IndexPage))
	router.GET("/read", Protect(LinkReadPage))
	router.GET("/unread", Protect(LinkUnreadPage))
	router.GET("/link/add", Protect(LinkAddPage))
	router.GET("/link/edit/:id", Protect(LinkEditPage))
	router.GET("/link/delete/:id", Protect(LinkDel))
	router.GET("/tag/delete/:name", Protect(TagDel))

	router.GET("/tags", Protect(TagsPage))
	router.GET("/static/:assettype/:filename", AssetsFinder)
}
