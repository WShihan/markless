package page

import (
	"markee/hooks"
	"markee/injection"
	"markee/model"
)

var (
	Env injection.Env
)

func LoadPage(router *model.RouterWithPrefix, env injection.Env) {
	Env = env
	router.GET("/login", LoginPage)
	router.GET("/setting", hooks.Protect(SettingPage))
	router.GET("/", hooks.Protect(LinkAllPage))
	router.GET("/read", hooks.Protect(LinkReadPage))
	router.GET("/unread", hooks.Protect(LinkUnreadPage))
	router.GET("/link/add", hooks.Protect(LinkAddPage))
	router.GET("/link/edit/:id", hooks.Protect(LinkEditPage))

	router.GET("/tags", hooks.Protect(TagsPage))
	router.GET("/static/:assettype/:filename", AssetsFinder)
}
