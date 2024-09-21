package handler

import (
	"github.com/julienschmidt/httprouter"
)

// 自定义路由器类型，包含一个前缀
type RouterWithPrefix struct {
	BaseURL string // 根路由
	Mux     *httprouter.Router
}

// 带前缀的 GET 方法
func (r *RouterWithPrefix) GET(path string, handle httprouter.Handle) {
	r.Mux.GET(r.BaseURL+path, handle)
}

// 带前缀的 POST 方法
func (r *RouterWithPrefix) POST(path string, handle httprouter.Handle) {
	r.Mux.POST(r.BaseURL+path, handle)
}

// 带前缀的 PUT 方法
func (r *RouterWithPrefix) PUT(path string, handle httprouter.Handle) {
	r.Mux.PUT(r.BaseURL+path, handle)
}

// 带前缀的 PATCH 方法
func (r *RouterWithPrefix) PATCH(path string, handle httprouter.Handle) {
	r.Mux.PATCH(r.BaseURL+path, handle)
}

// 带前缀的 DELETE 方法
func (r *RouterWithPrefix) DELETE(path string, handle httprouter.Handle) {
	r.Mux.DELETE(r.BaseURL+path, handle)
}

// 带前缀的 OPTIONS 方法
func (r *RouterWithPrefix) OPTIONS(path string, handle httprouter.Handle) {
	r.Mux.OPTIONS(r.BaseURL+path, handle)
}

// 带前缀的 OPTIONS 方法
// func (r *RouterWithPrefix) NotFound(path string, handle httprouter.Handle) {
// 	r.Mux.NotFound = http.HandlerFunc(handle)
// }
