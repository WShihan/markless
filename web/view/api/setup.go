package api

import (
	"markless/injection"
	"markless/web/handler"
)

var (
	Env injection.Env
)

func LoadAPI(router *handler.RouterWithPrefix) {
	router.POST("/api/link/add", handler.Protect(LinkAddApi))
	router.POST("/api/tag/update/name", handler.Protect(TagUpdateName))
	router.POST("/api/tag/update/link", handler.Protect(TagUpdateLink))

	router.GET("/api/link/all", handler.Protect(LinkAllApi))
	router.GET("/api/link/delete/:id", handler.Protect(LinkDel))
	router.GET("/api/link/read/:id", handler.Protect(LinkRead))
	router.GET("/api/link/unread/:id", handler.Protect(LinkUnread))

	router.GET("/api/tag/delete/:name", handler.Protect(TagDelApi))
	// router.GET("/api/tag/delete/:name", handler.Protect(TagDelApi))

}
