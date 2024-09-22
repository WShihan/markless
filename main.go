package main

import (
	"flag"
	"markless/injection"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/handler"
	"markless/web/view/api"
	"markless/web/view/page"
)

var (
	Commit    string
	Version   string
	BuildTime string
)

func main() {
	Title := flag.String("title", "markless", "网站名称")
	BaseURL := flag.String("baseurl", "", "根路由")
	DataBaseURL := flag.String("databaseurl", tool.ExcutePath()+"/markless.db", "数据库地址")
	Port := flag.Int("port", 5000, "运行端口")
	adminName := flag.String("adminname", "admin", "初始用户名称")
	adminPassword := flag.String("password", "admin1234", "初始用户密码")
	flag.Parse()

	env := injection.Env{
		BaseURL:     *BaseURL,
		Title:       *Title,
		DataBaseURL: *DataBaseURL,
		Port:        *Port,
		Version:     Version,
		Commit:      Commit,
		BuildTime:   BuildTime,
	}
	util.InitENV(&env)
	store.InitDB(*DataBaseURL)
	store.InitAdmin(*adminName, *adminPassword)
	handler.InitEnv(&env)
	api.LoadAPI(&handler.Router)
	page.LoadPage(&handler.Router, &env)
	handler.RunServer()
}
