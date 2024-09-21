package util

import (
	"html/template"
	"markless/local"
	"markless/tool"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"JoinTagNames": tool.JoinTagNames,
		"Increase":     tool.Increase,
		"Decrease":     tool.Decrease,
		"TimeFMT":      tool.TimeFMT,
		"i18n":         local.Translate,
	}
}

func GetBaseTemplate() *template.Template {
	funcMap := GetFuncMap()
	return template.New("html/template.html").Funcs(funcMap)
}
