package util

import (
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
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
