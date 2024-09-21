package local

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestTranslater(t *testing.T) {
	localizer := i18n.NewLocalizer(Bundle, "zh-CN")
	// 获取翻译
	hello, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "Hello"})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%v", hello)
}

func TestUnfind(t *testing.T) {
	k := "nav.setting"
	lang := "zh-CN"
	s := Translate(k, lang)
	t.Logf("%s", s)
}
