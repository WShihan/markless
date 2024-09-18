package tool

import (
	"fmt"
	"html/template"
	"log/slog"
	"markee/model"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FileOrPathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func GetBaseTemplate() *template.Template {
	return template.New("templates/temp.html")
}
func ExcutePath() string {
	excutePath, err := os.Executable()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return filepath.Dir(excutePath)
}

func Find(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// 自定义函数，用于合并多个 User 对象的 FirstName 字段
func JoinTagNames(tags []model.Tag) string {
	firstNames := make([]string, len(tags))
	for i, user := range tags {
		firstNames[i] = user.Name
	}
	return strings.Join(firstNames, ", ")
}

// 自增模板函数
func Increase(num int) int {
	return num + 1
}

// 自增模板函数
func Decrease(num int) int {
	return num - 1
}

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
func TimeFMT(t time.Time) string {
	day := t.Day()
	month := t.Month()
	year := t.Year()
	hour := t.Hour()
	min := t.Minute()

	return fmt.Sprintf("%d-%d-%d %d:%d", year, month, day, hour, min)
}
