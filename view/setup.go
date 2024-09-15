package view

import (
	"marky/model"
)

var (
	BASE_URL = "/marky"
)

func LoadApi(router *model.RouterWithPrefix) {
	router.POST("/user/login", UserLogin)
	router.POST("/link/create", LinkCreate)

}

func LoadView(router *model.RouterWithPrefix) {
	router.GET("/", Protect(IndexPage))
	router.GET("/login", Login)
	router.GET("/link/add", Protect(LinkAddPage))
	router.GET("/link/delete/:id", Protect(LinkDel))
	router.GET("/static/:assettype/:filename", AssetsFinder)
}
