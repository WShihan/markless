package main

import (
	"flag"
	"fmt"
	"markee/logging"
	"markee/model"
	"markee/store"
	"markee/util"
	"markee/view"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	BaseURL := flag.String("baseurl", "", "根路由")
	DataBaseURL := flag.String("databaseurl", util.ExcutePath()+"/markee.db", "数据库地址")
	Port := flag.Int("port", 5000, "运行端口")
	adminname := flag.String("adminname", "wsh", "初始用户名称")
	adminPassword := flag.String("adminpassword", "test123", "初始用户密码")

	flag.Parse()
	logging.InitLogger()

	Mux := *httprouter.New()
	// 创建自定义路由器，指定前缀为 /app/go
	router := &model.RouterWithPrefix{
		BaseURL: *BaseURL,
		Router:  &Mux,
	}
	store.InitDB(*DataBaseURL)
	env := model.BaseInjdection{
		BaseURL: *BaseURL,
		Title:   "markee",
	}
	view.LoadApi(router)
	view.LoadPage(router, env)
	view.InitAdmin(*adminname, *adminPassword)
	runAt := fmt.Sprintf("127.0.0.1:%d", *Port)

	cos := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Custom-Header", "uid"},
		AllowCredentials: true,
	})
	handler := cos.Handler(view.ApplyMiddleware(&Mux, view.LogRequest))
	server := http.Server{
		Addr:    runAt,
		Handler: handler,
	}
	logging.Logger.Info("server run in:", runAt+*BaseURL)
	server.ListenAndServe()
}
