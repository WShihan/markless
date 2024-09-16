package util

import (
	"html/template"
	"log/slog"
	"markee/model"
	"os"
	"path/filepath"
	"strings"
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

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"JoinTagNames": JoinTagNames,
	}
}
