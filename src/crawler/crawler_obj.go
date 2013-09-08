package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"time"
)

type Crawler struct {
	last_crawl time.Time
	website    string
	nmangas    int
	base_path  string
	semaphores chan bool
	mangas     map[string]Manga
}

func (c *Crawler) manga_list() []string {
	list := []string{}
	var doc *goquery.Document
	var link string
	switch c.website {
	case "mangareader":
		doc = goquery_download_safe("http://www.mangareader.net/alphabetical")
		doc.Find("ul.series_alpha a").Each(func(i int, s *goquery.Selection) {
			link, _ = s.Attr("href")
			list = append(list, link)
		})
	default:
		fmt.Printf(">website not implemented %s", c.website)
	}
	return list
}

func (c *Crawler) create_manga_entry(link string) Manga {
	var rl int64
	var m Manga
	//	fmt.Println("Creating")
	doc := goquery_download_safe(fmt.Sprintf("%s%s", c.base_path, link))
	tags := []string{}
	m.name = doc.Find("h2.aname").Text()
	m.rel_path = link
	doc.Find("div#mangaproperties table td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 3:
			m.alt_names = s.Text()
		case 5:
			rl, _ = strconv.ParseInt(s.Text(), 10, 0)
			m.release_year = int(rl)
		case 7:
			m.ongoing = (s.Text() != "Complete")
		case 9:
			m.author = s.Text()
		case 11:
			m.artist = s.Text()
		case 13:
			m.normal = (s.Text() == "Right to Left")
		}
	})
	m.parent = c
	doc.Find("span.genretags").Each(func(i int, s *goquery.Selection) {
		tags = append(m.tags, s.Text())
	})
	m.tags = tags
	c.nmangas--
	return m
}

func (c *Crawler) update_mangas() {
	mlist := c.manga_list()
	c.nmangas = cap(mlist)
	for _, link := range mlist {
		if _, present := c.mangas[link]; !present {
			go func(c *Crawler, link string, sem chan bool) {
				<-sem
				c.mangas[link] = c.create_manga_entry(link)
				sem <- true
			}(c, link, c.semaphores)
		}
	}
	//Test stuff
}

func (c *Crawler) initialize_semaphors(n int) {
	c.semaphores = make(chan bool, n)
	for i := 0; i < n; i++ {
		c.semaphores <- true
	}
}
