package page

import (
	"markless/injection"
	"markless/store"
	"markless/tool"
	"markless/util"
	"markless/web/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssetsFinder(t *testing.T) {
	env := injection.Env{
		BaseURL:     "/webapp/markless",
		DataBaseURL: tool.ExcutePath() + "/markless.db",
		Title:       "markless",
		Port:        5000,
	}
	util.InitENV(&env)
	util.InitLogger()
	store.InitDB(env.DataBaseURL)
	store.InitAdmin("admin", "admin1234")
	handler.InitEnv(&env)
	LoadPage(&handler.Router, &env)

	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", env.BaseURL+"/static/css/test.css", nil)
	if err != nil {
		t.Fatal("错误：", err)
	}

	rr := httptest.NewRecorder()
	handler.Router.Mux.ServeHTTP(rr, req)

	// 检查响应状态码
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	t.Logf("assets content:%v",
		rr.Body.String())

}
