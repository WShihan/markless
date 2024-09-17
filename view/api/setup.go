package api

import (
	"markee/hooks"
	"markee/injection"
	"markee/model"
)

var (
	Env injection.Env
)

func LoadAPI(router *model.RouterWithPrefix) {
	router.POST("/user/login", UserLogin)
	router.POST("/user/password", hooks.Protect(UserChangePassword))
	router.POST("/user/token/add", hooks.Protect(UserTokenAdd))
	router.POST("/user/token/delete", hooks.Protect(UserTokenDelete))
	router.POST("/link/add", hooks.Protect(LinkAdd))
	router.POST("/link/update/:id", hooks.Protect(LinkUpdate))
	router.GET("/link/delete/:id", hooks.Protect(LinkDel))
	router.GET("/link/read/:id", hooks.Protect(LinkRead))
	router.GET("/link/unread/:id", hooks.Protect(LinkUnread))

	router.GET("/tag/delete/:name", hooks.Protect(TagDel))
	router.POST("/tag/add", hooks.Protect(TagAdd))
}
