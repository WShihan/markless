package api

import (
	"bytes"
	"encoding/json"
	"markless/injection"
	"markless/model"
	"markless/store"
	"markless/tool"
	"markless/util"
	handler "markless/web/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLinkApi(t *testing.T) {
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
	LoadAPI(&handler.Router)

	// 创建书签
	postdata := LinkAddPost{
		Url:  "https://www.baidu.com",
		Desc: "百度",
		Tags: "标签1&标签2",
		Read: true,
	}
	jsonData, _ := json.Marshal(postdata)

	req, err := http.NewRequest(http.MethodPost, env.BaseURL+"/api/link/add", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal("错误：", err)
	}
	user := model.User{}
	store.DB.Where("username = ?", "admin").First(&user)
	t.Logf("toeken:%v", *user.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", *user.Token)

	rr := httptest.NewRecorder()
	handler.Router.Mux.ServeHTTP(rr, req)

	// 检查响应状态
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v", rr.Body.String())
	}

	t.Logf("content:%v", rr.Body.String())

}

func TestLinkAll(t *testing.T) {
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
	LoadAPI(&handler.Router)

	// 获取所有书签
	req, err := http.NewRequest(http.MethodGet, env.BaseURL+"/api/link/all", nil)
	if err != nil {
		t.Fatal("错误：", err)
	}
	user := model.User{}
	store.DB.Where("username = ?", "admin").First(&user)
	t.Logf("toeken:%v", *user.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", *user.Token)

	rr := httptest.NewRecorder()
	handler.Router.Mux.ServeHTTP(rr, req)

	// 检查响应状态
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v", rr.Body.String())
	}

	t.Logf("content:%v", rr.Body.String())

}

func TestLinkDelApi(t *testing.T) {
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
	LoadAPI(&handler.Router)
	// 删除元素

	req, err := http.NewRequest(http.MethodGet, env.BaseURL+"/api/link/all", nil)
	if err != nil {
		t.Fatal("错误：", err)
	}
	user := model.User{}
	store.DB.Where("username = ?", "admin").First(&user)
	t.Logf("toeken:%v", *user.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", *user.Token)

	rr := httptest.NewRecorder()
	handler.Router.Mux.ServeHTTP(rr, req)

	// 检查响应状态
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v", rr.Body.String())
	}

	t.Logf("content:%v", rr.Body.String())

}
