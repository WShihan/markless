package main

import (
	"fmt"
	"log"
	"time"

	readability "github.com/go-shiori/go-readability"
)

var (
	urls = []string{
		"http://www.zhexueshi.com/paper/9210",
		"https://gorm.io/zh_CN/docs/associations.html#%E5%88%A0%E9%99%A4%E5%85%B3%E8%81%94",
	}
)

func parse(url string) (map[string]string, error) {
	infoMap := make(map[string]string)
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return infoMap, fmt.Errorf("failed to parse %s, %v", url, err)
	}
	infoMap["url"] = url
	infoMap["title"] = article.Title
	infoMap["author"] = article.Byline
	infoMap["desc"] = article.Excerpt
	infoMap["sitename"] = article.SiteName
	infoMap["image"] = article.Image
	infoMap["icon"] = article.Favicon
	infoMap["content"] = article.Content

	return infoMap, nil
}

func main() {
	for _, url := range urls {
		infoMap, _ := parse(url)
		log.Println(infoMap)
	}
}
