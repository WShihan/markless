package local

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	Bundle *i18n.Bundle
)

const defaultLang = "zh-CN"

//go:embed translations/*.json
var TranslationsFS embed.FS

func init() {
	Bundle = i18n.NewBundle(language.Chinese)

	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	Bundle.LoadMessageFileFS(TranslationsFS, "translations/zh-CN.json")
	Bundle.LoadMessageFileFS(TranslationsFS, "translations/zh-TW.json")
	Bundle.LoadMessageFileFS(TranslationsFS, "translations/en.json")
	Bundle.LoadMessageFileFS(TranslationsFS, "translations/ja.json")

}

func Translate(k string, lang string) string {
	tag := language.Make(lang).String()
	localizer := i18n.NewLocalizer(Bundle, tag)
	translated, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: k})
	if err != nil {
		defaultlocalizer := i18n.NewLocalizer(Bundle, defaultLang)
		translated, err := defaultlocalizer.Localize(&i18n.LocalizeConfig{MessageID: k})
		if err != nil {
			return k
		}
		return translated
	}

	return translated

}

func GetPreferredLanguage(r *http.Request) string {
	acceptLang := r.Header.Get("Accept-Language")
	// reqeust headers have not Accept-Language
	if acceptLang == "" {
		return defaultLang
	}

	// parse Accept-Language
	languages := strings.Split(acceptLang, ",")
	if len(languages) > 0 {
		return strings.TrimSpace(languages[0])
	}

	return defaultLang
}
