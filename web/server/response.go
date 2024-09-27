package server

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Status bool        `json:"status"`
}

func ApiSuccess(w *http.ResponseWriter, data interface{}) {
	(*w).Header().Set("Content-Type", "application/json")
	res := &ApiResponse{
		Status: true,
		Data:   data,
	}
	jsonData, _ := json.Marshal(res)
	(*w).Write(jsonData)
}

func ApiFailed(w *http.ResponseWriter, code int, msg string) {
	(*w).Header().Set("Content-Type", "application/json")
	res := &ApiResponse{
		Msg:    msg,
		Status: false,
	}
	jsonData, _ := json.Marshal(res)
	(*w).Write(jsonData)
}
