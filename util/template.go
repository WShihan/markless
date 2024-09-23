package util

import (
	"html/template"
	"markless/local"
	"markless/tool"
)

func RenderHTML(HTMLString string) template.HTML {
	return template.HTML(HTMLString)
}

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"JoinTagNames": tool.JoinTagNames,
		"Increase":     tool.Increase,
		"Decrease":     tool.Decrease,
		"TimeFMT":      tool.TimeFMT,
		"RandomN":      tool.RandomN,
		"i18n":         local.Translate,
		"HTML":         RenderHTML,
	}
}

func GetBaseTemplate() *template.Template {
	funcMap := GetFuncMap()
	return template.New("html/template.html").Funcs(funcMap)
}
