package util

import (
	"testing"
)

func TestFunc(t *testing.T) {
	pageinfo, err := Scrape("https://www.baidu.com", 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%v", pageinfo.Preview)
}
