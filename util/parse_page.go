package util

import (
	"fmt"
	"time"

	readability "github.com/go-shiori/go-readability"
)

type PageInfo struct {
	Url      string
	Title    string
	Author   string
	SiteName string
	Image    string
	Icon     string
	Content  string
	Desc     string
}

func ParsePage(url string) (*PageInfo, error) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s, %v", url, err)
	}
	pageInfo := &PageInfo{
		Url:      url,
		Title:    article.Title,
		Author:   article.Byline,
		SiteName: article.SiteName,
		Image:    article.Image,
		Icon:     article.Favicon,
		Content:  article.Content,
		Desc:     article.Excerpt,
	}

	return pageInfo, nil
}
