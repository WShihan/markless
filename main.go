package main

import (
	"flag"
	"fmt"
	"markee/hooks"
	"markee/injection"
	"markee/logging"
	"markee/model"
	"markee/store"
	"markee/tool"
	"markee/util"
	"markee/view/api"
	"markee/view/page"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	BaseURL := flag.String("baseurl", "", "根路由")
	DataBaseURL := flag.String("databaseurl", tool.ExcutePath()+"/markee.db", "数据库地址")
	Port := flag.Int("port", 5000, "运行端口")
	adminname := flag.String("adminname", "admin", "初始用户名称")
	adminPassword := flag.String("adminpassword", "markee1234", "初始用户密码")

	flag.Parse()
	logging.InitLogger()

	Mux := *httprouter.New()
	// 创建自定义路由器，指定前缀为 /app/go
	router := &model.RouterWithPrefix{
		BaseURL: *BaseURL,
		Router:  &Mux,
	}
	env := injection.Env{
		BaseURL: *BaseURL,
		Title:   "markee",
	}
	store.InitDB(*DataBaseURL)
	util.InitENV(env)
	util.InitAdmin(*adminname, *adminPassword)

	api.LoadAPI(router)
	page.LoadPage(router, env)
	runAt := fmt.Sprintf("127.0.0.1:%d", *Port)

	cos := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Custom-Header", "uid"},
		AllowCredentials: true,
	})
	handler := cos.Handler(hooks.ApplyMiddleware(&Mux, hooks.LogRequest))
	server := http.Server{
		Addr:    runAt,
		Handler: handler,
	}
	logging.Logger.Info("server run in:", runAt+*BaseURL)
	server.ListenAndServe()
}
