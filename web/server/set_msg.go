package server

import (
	"net/http"
	"net/url"
)

func SetMsg(w *http.ResponseWriter, message string) {
	http.SetCookie(*w, &http.Cookie{
		Name:  "message",
		Value: url.QueryEscape(message),
		Path:  "/",
	})
	http.SetCookie(*w, &http.Cookie{
		Name:  "message_shown",
		Value: "false",
		Path:  "/",
	})
}
