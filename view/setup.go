package view

import (
	"log/slog"
	"marky/model"
	"net/http"
)

var (
	BASE_URL = "/marky"
	Env      model.BaseInjdection
)

func Redirect(w http.ResponseWriter, r *http.Request, route string) {
	http.Redirect(w, r, Env.BaseURL+route, http.StatusMovedPermanently)
	slog.Info("redirest:" + Env.BaseURL + route)
}

func LoadApi(router *model.RouterWithPrefix) {
	router.POST("/user/login", UserLogin)
	router.POST("/link/add", Protect(LinkAdd))
	router.POST("/link/update/:id", Protect(LinkUpdate))
	router.POST("/tag/add", Protect(TagAdd))

}

func LoadPage(router *model.RouterWithPrefix, env model.BaseInjdection) {
	Env = env
	router.GET("", Protect(IndexPage))
	router.GET("/login", Login)
	router.GET("/link/add", Protect(LinkAddPage))
	router.GET("/link/edit/:id", Protect(LinkEditPage))
	router.GET("/link/delete/:id", Protect(LinkDel))
	router.GET("/link/read/:id", LinkRead)
	router.GET("/tag/delete/:name", Protect(TagDel))

	router.GET("/tags", Protect(TagsPage))
	router.GET("/static/:assettype/:filename", AssetsFinder)
}
