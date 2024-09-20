package page

import (
	"markless/injection"
	"markless/web/handler"
)

var (
	Env injection.Env
)

func LoadPage(router *handler.RouterWithPrefix, env *injection.Env) {
	Env = *env
	router.GET("/login", LoginPage)
	router.GET("/", handler.Protect(IndexPage))
	router.GET("/all", handler.Protect(LinkAllPage))
	router.GET("/read", handler.Protect(LinkReadPage))
	router.GET("/unread", handler.Protect(LinkUnreadPage))
	router.GET("/link/find", handler.Protect(LinkAddPage))
	router.GET("/link/edit/:id", handler.Protect(LinkEditPage))
	router.GET("/tags", handler.Protect(TagsPage))
	router.GET("/tag/edit/:id", handler.Protect(TagEditPage))
	router.GET("/register", RegisterPage)
	router.GET("/setting", handler.Protect(SettingPage))
	router.GET("/static/:assettype/:filename", AssetsFinder)

	router.POST("/link/add", handler.Protect(LinkAdd))
	router.POST("/link/update/:id", handler.Protect(LinkUpdate))
	router.POST("/tag/add", handler.Protect(TagAdd))
	router.POST("/user/login", UserLogin)
	router.POST("/user/register", UserRegister)
	router.POST("/user/password", handler.Protect(UserChangePassword))
	router.POST("/user/token/add", handler.Protect(UserTokenAdd))
	router.POST("/user/token/delete", handler.Protect(UserTokenDelete))
}
