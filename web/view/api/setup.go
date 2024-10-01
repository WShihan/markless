package api

import (
	"markless/injection"
	handler "markless/web/server"
)

var (
	Env injection.Env
)

func LoadAPI(router *handler.RouterWithPrefix, env *injection.Env) {
	Env = *env
	router.POST("/api/link/add", handler.Protect(LinkAdd))
	router.POST("/api/link/markread", handler.Protect(MarkAllAsReadOrRead))
	router.POST("/api/link/exist", handler.Protect(LinkExist))
	router.POST("/api/link/update", handler.Protect(LinkUpdate))
	router.POST("/api/link/attach", handler.Protect(LinkAttachTags))
	router.POST("/api/tag/update/name", handler.Protect(TagUpdateName))
	router.POST("/api/tag/update/link", handler.Protect(TagUpdateLink))
	router.POST("/api/tag/add", handler.Protect(TagAdd))
	router.POST("/api/tag/attach", handler.Protect(AttachLinks))
	router.POST("/api/user/login", UserLogin)
	router.POST("/api/user/register", UserRegister)
	router.POST("/api/user/info/update", handler.Protect(UserInfoUpdate))
	router.POST("/api/user/password/update", handler.Protect(UserUpdatePassword))

	router.GET("/api/link/all", handler.Protect(LinkAll))
	router.GET("/api/link/pagination", handler.Protect(LinkPagination))
	router.GET("/api/link/delete/:id", handler.Protect(LinkDel))
	router.GET("/api/link/read/:id", handler.Protect(LinkRead))
	router.GET("/api/link/unread/:id", handler.Protect(LinkUnread))
	router.GET("/api/tag/all", handler.Protect(TagAll))
	router.GET("/api/tag/delete/:name", handler.Protect(TagDelApi))
	router.GET("/api/tag/related-link/:name", handler.Protect(TagRelatedLinks))
	router.GET("/api/tag/stastic", handler.Protect(TagStastic))

	router.GET("/api/link/one/:id", handler.Protect(LinkOne))
	router.GET("/api/link/archive/update/:id", handler.Protect(LinkUpdateArchive))
	router.GET("/api/user/info", handler.Protect(UserInfoGet))
	router.GET("/api/user/env", handler.Protect(UserEnvGet))
	router.GET("/api/user/token/delete", handler.Protect(UserTokenDelete))
	router.GET("/api/user/token/refresh", handler.Protect(UserTokenRefresh))

}
