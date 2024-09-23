package server

import (
	"fmt"
	"markless/injection"
	"markless/util"
	"net/http"

	_ "net/http/pprof"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var (
	Env    *injection.Env
	Router RouterWithPrefix
	Server http.Server
)

func InitEnv(env *injection.Env) {
	Env = env
	mux := httprouter.New()
	Router = RouterWithPrefix{
		Mux:     mux,
		BaseURL: Env.BaseURL,
	}
	cos := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Custom-Header", "uid", "X-Token"},
		AllowCredentials: true,
	})
	handler := cos.Handler(ApplyHooks(Router.Mux, LogRequest))
	runAt := fmt.Sprintf("127.0.0.1:%d", Env.Port)
	Server = http.Server{
		Addr:    runAt,
		Handler: handler,
	}
}

func RunServer(env *injection.Env) {
	util.Logger.Info(fmt.Sprintf("version:%s\tcommit:%s\tbuild-time:%s", env.Version, env.Commit, env.BuildTime))
	util.Logger.Info(fmt.Sprintf("server run in:\thttp://%s", fmt.Sprintf("127.0.0.1:%d", Env.Port)+Env.BaseURL))
	util.Logger.Info(
		`
							
							■■■■■■■■■            
							■■■■■■■■■            
							■■■■■■■■■            
							■■■■■■■■■            
							■■■■■■■■■            
							■■■■ ■■■■            
							■■■   ■■■            
							■■     ■■            
							■       ■            
									  
						    Welcome to Markless            
							
		`)
	err := Server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
