package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"time"
	"strconv"
)

type Chapter struct {
	number int
	pages int
	title string
}

type Manga struct{
	name string
	rel_path string
	alt_names string
	release_year int
	author string
	artist string
	ongoing bool
	normal bool //reading direction, normal <-
	tags []string
	last_chapter int
	abstract string
}

func (m *Manga) chapter_list_mangareader(base_path string) []string {
	chap_list := []string{};
	var link string
	doc,_ := goquery.NewDocument(fmt.Sprintf("%s%s",base_path,m.rel_path))
	doc.Find("table#listing a").Each(func(i int, s *goquery.Selection){
		link,_ = s.Attr("href")
		chap_list= append(chap_list,link)
	})
	return chap_list
}

type Crawler struct {
	last_crawl time.Time
	website string
	nmangas int
	base_path string
	semaphores chan bool
	mangas map[string]Manga
}

func (c *Crawler) manga_list() []string {
	list := []string{}
	var doc *goquery.Document
	var link string
	switch c.website{
	case "mangareader":
		doc,_ = goquery.NewDocument("http://www.mangareader.net/alphabetical")
		doc.Find("ul.series_alpha a").Each(func(i int, s *goquery.Selection) {
			link,_ = s.Attr("href")
			list= append(list,link)
		})
	default:
		fmt.Printf(">website not implemented %s",c.website)
	}
	return list
}

func (c *Crawler) create_manga_entry(link string) Manga{
	var rl int64
	var m Manga
//	fmt.Println("Creating")
	doc,err := goquery.NewDocument(fmt.Sprintf("%s%s",c.base_path,link))
	if err != nil {
		fmt.Println("Found an Error!")
		fmt.Println(err)
	}
	tags := []string{}
	m.name = doc.Find("h2.aname").Text()
	doc.Find("div#mangaproperties table td").Each(func(i int,s *goquery.Selection){
		switch i{
		case 3:
			m.alt_names = s.Text()
		case 5:
			rl,_ = strconv.ParseInt(s.Text(),10,0)
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
	doc.Find("span.genretags").Each(func(i int,s *goquery.Selection){
		tags= append(m.tags, s.Text())
	})
	m.tags= tags
	fmt.Printf("%s\n",m.name)
	c.nmangas--
	return m
}

func (c *Crawler) update_mangas() {
	mlist := c.manga_list()
	c.nmangas = cap(mlist)
	for _,link := range mlist {
		if _,present := c.mangas[link]; !present{
			go func(c *Crawler, link string, sem chan bool){
				<-sem
				c.mangas[link]= c.create_manga_entry(link)
				sem <- true
			}(c,link, c.semaphores)
		}
	}
	//Test stuff
	fmt.Println("Done")
}

func main() {
	crawl := Crawler{time.Now(),"mangareader",0,"http://mangareader.net",make(chan bool,100),make(map[string]Manga)}
	for i:=0;i<20;i++{
		crawl.semaphores<-true
	}
	crawl.update_mangas()
	for crawl.nmangas > 0{
		time.Sleep(time.Second)
	}
	for _,manga := range crawl.mangas {
		fmt.Printf("%s\n",manga.name)
	}
}
