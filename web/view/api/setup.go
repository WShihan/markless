package api

import (
	"markless/injection"
	handler "markless/web/server"
)

var (
	Env injection.Env
)

func LoadAPI(router *handler.RouterWithPrefix) {
	router.POST("/api/link/add", handler.Protect(LinkAdd))
	router.POST("/api/tag/update/name", handler.Protect(TagUpdateName))
	router.POST("/api/tag/update/link", handler.Protect(TagUpdateLink))
	router.POST("/api/link/exist", handler.Protect(LinkExist))

	router.GET("/api/link/all", handler.Protect(LinkAll))
	router.GET("/api/link/delete/:id", handler.Protect(LinkDel))
	router.GET("/api/link/read/:id", handler.Protect(LinkRead))
	router.GET("/api/link/unread/:id", handler.Protect(LinkUnread))
	router.GET("/api/tag/delete/:name", handler.Protect(TagDelApi))

}
