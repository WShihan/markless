package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// 获取网页内容
func fetchURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// 提取网页标题
func getTitle(doc *html.Node) string {
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}

// 提取网页 favicon
func getFavicon(doc *html.Node, baseURL *url.URL) string {
	var favicon string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, href string
			for _, attr := range n.Attr {
				if attr.Key == "rel" {
					rel = attr.Val
				}
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			if strings.Contains(rel, "icon") {
				favicon = resolveURL(href, baseURL)
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return favicon
}

// 解析相对 URL
func resolveURL(href string, baseURL *url.URL) string {
	parsedURL, err := url.Parse(href)
	if err != nil {
		return href
	}
	if parsedURL.IsAbs() {
		return href
	}
	return baseURL.ResolveReference(parsedURL).String()
}

func main() {
	rawURL := "https://www.wsh233.cn/blog/88525144"
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	doc, err := fetchURL(rawURL)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	title := getTitle(doc)
	favicon := getFavicon(doc, baseURL)

	fmt.Println("Title:", title)
	fmt.Println("Favicon:", favicon)
}
