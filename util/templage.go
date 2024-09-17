package util

import (
	"html/template"
	"markee/tool"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"JoinTagNames": tool.JoinTagNames,
		"Increase":     tool.Increase,
		"Decrease":     tool.Decrease,
	}
}

func GetBaseTemplate() *template.Template {
	funcMap := GetFuncMap()
	return template.New("html/template.html").Funcs(funcMap)
}
