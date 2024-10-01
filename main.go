package main

import (
	"flag"
	"markless/injection"
	"markless/store"
	"markless/tool"
	"markless/util"
	handler "markless/web/server"
	"markless/web/view/api"
	"markless/web/view/page"
)

var (
	Commit    string
	Version   string
	BuildTime string
)

func main() {
	Title := flag.String("title", "markless", "App Name")
	BaseURL := flag.String("baseurl", "", "Base UTL")
	JWTExpire := flag.Int("jwtexpire", 60, "JWT expires time in minutes")
	DataBaseURL := flag.String("databaseurl", tool.ExcutePath()+"/markless.db", "Path to database file")
	Port := flag.Int("port", 5000, "Port")
	adminName := flag.String("adminname", "admin", "Iitial administrator user name")
	adminPassword := flag.String("password", "admin1234", "Initial administrator user password")
	flag.Parse()

	env := injection.Env{
		BaseURL:     *BaseURL,
		Title:       *Title,
		DataBaseURL: *DataBaseURL,
		Port:        *Port,
		Version:     Version,
		Commit:      Commit,
		BuildTime:   BuildTime,
		HmacSecret:  tool.ShortUID(12),
		SecretKey:   tool.ShortUID(32),
		JWTExpire:   *JWTExpire,
	}
	util.InitENV(&env)
	store.InitDB(*DataBaseURL)
	store.InitAdmin(*adminName, *adminPassword)
	handler.InitEnv(&env)
	api.LoadAPI(&handler.Router, &env)
	page.LoadPage(&handler.Router, &env)
	handler.RunServer(&env)
}
