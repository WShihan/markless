package util

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func Parse_Webpage(url string) (string, string) {
	title, favicon := getTitleAndFavicon(url)
	return title, favicon
}

func getTitleAndFavicon(url string) (string, string) {
	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return "", ""
	}
	defer resp.Body.Close()

	// 解析 HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing the HTML:", err)
		return "", ""
	}

	var title, favicon string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, attr := range n.Attr {
				if attr.Key == "rel" && (attr.Val == "icon" || attr.Val == "shortcut icon") {
					// 获取 favicon 的地址
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							favicon = attr.Val
							break
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// 处理 favicon 地址，确保是完整的 URL
	if favicon != "" && !strings.HasPrefix(favicon, "http") {
		favicon = url + favicon // 这里简单处理，可能需要更复杂的逻辑
	}

	return title, favicon
}
