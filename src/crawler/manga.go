package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

type Manga struct {
	name         string
	rel_path     string
	alt_names    string
	release_year int
	author       string
	artist       string
	ongoing      bool
	normal       bool //reading direction, normal <-
	tags         []string
	abstract     string
	last_chapter int
	chapters     []*Chapter
	parent       *Crawler
}

func (m *Manga) chapter_list_mangareader() {
	// Chapter title parsing should be done here.
	m.chapters = *new([]*Chapter)
	var tchap *Chapter
	var link string
	doc := goquery_download_safe(fmt.Sprintf("%s%s", m.parent.base_path, m.rel_path))
	doc.Find("table#listing a").Each(func(i int, s *goquery.Selection) {
		link, _ = s.Attr("href")
		tchap = &Chapter{
			title:         s.Text(),
			number:        i,
			starting_path: fmt.Sprintf("%s%s", m.parent.base_path, link),
			done:          make(chan bool),
			parent:        m.parent,
		}
		//FIX: format this later
		m.chapters = append(m.chapters, tchap)
	})
}

func (m *Manga) download_chapter_list() {
	// Calls chapter_list_$, with the links to the first page of each chapter,
	// calls Chapter.download_chapter(#) on each string of the returned slice
	// For now this only works with mangareader
	var tmpstring, link string
	var doc *goquery.Document
	for _, chap := range m.chapters {
		doc = goquery_download_safe(chap.starting_path)
		tmpstring = doc.Find("div#selectpage").Text()
		chap.pages, _ = strconv.Atoi(tmpstring[len(tmpstring)-2:])
		link, _ = doc.Find("img#img").Attr("src")
		go chap.download_chapter(link)
	}
	for _, chap := range m.chapters {
		<-chap.done
	}
}
