package page

import (
	"markless/injection"
	handler "markless/web/server"
)

var (
	Env injection.Env
)

func LoadPage(router *handler.RouterWithPrefix, env *injection.Env) {
	Env = *env
	router.GET("/", IndexPage)
	router.GET("/static/:assettype/:filename", AssetsFinder)
}
