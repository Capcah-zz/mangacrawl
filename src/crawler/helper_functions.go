package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func http_download_safe(path string) *http.Response {
	content, err := http.Get(path)
	for err != nil {
		content, err = http.Get(path)
	}
	return content
}

func goquery_download_safe(path string) *goquery.Document {
	doc, err := goquery.NewDocument(path)
	for err != nil {
		doc, err = goquery.NewDocument(path)
	}
	return doc
}
