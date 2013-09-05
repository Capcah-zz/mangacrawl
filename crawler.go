package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"time"
	"strconv"
	"math/rand"
	"net/http"
	"io"
	"archive/tar"
)

type Chapter struct {
	number int
	pages int
	title string
	starting_path string
	done chan bool
	parent * Crawler
}

func (chap *Chapter) download_page(path string, ret_chan *chan io.ReadCloser) {
	<-chap.parent.semaphores
	content,_ := http.Get(path)
	*ret_chan <- content.Body
	chap.parent.semaphores<-true
}

func (chap *Chapter) store_chapter(pages_data []io.ReadCloser) {
	var header *tar.Header
	for i,data:= range pages_data {
		header = new(tar.Header)
		header.Name = strconv.Itoa(i)
		header.Size = (1<<30)
		header.ModTime = time.Now()
		header.Mode = 0644
		chap.tar_writer.WriteHeader(header)
		io.Copy(chap.tar_writer,data)
	}
}

func (chap *Chapter) download_chapter()  {
	fd,_ := os.Create(chap.title)
	defer fd.Close()
	gw := gzip.NewWriter(fd)
	defer gw.Close()
	chap.tar_writer = tar.NewWriter(gw)
	defer chap.tar_writer.Close()
	sweirdnum,eweirdnum := strings.Index(chap.starting_path,"-"),strings.LastIndex(
								chap.starting_path,".")
	weirdnum,_ := strconv.Atoi(chap.starting_path[sweirdnum+1:eweirdnum])
	base_dl_path := chap.starting_path[0:sweirdnum+1]
	pages_chanel := make([]chan io.ReadCloser,chap.pages)
	pages_data := make([]io.ReadCloser,chap.pages)
	for i,_ := range pages_chanel {
		pages_chanel[i] = make(chan io.ReadCloser)
	}
	for i,_ := range pages_data{
		fmt.Println(i)
		go chap.download_page(fmt.Sprintf("%s%d.jpg",base_dl_path ,weirdnum+2*i),
									&pages_chanel[i])
	} 
	for i,_ := range pages_data{
		fmt.Println(">",i)
		pages_data[i] = <-(pages_chanel[i])
	}
	chap.store_chapter(pages_data)
	chap.done <- true
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
	abstract string
	last_chapter int
	chapters []Chapter
	parent * Crawler
}

func (m *Manga) chapter_list_mangareader() []string {
	// This function only defines number, 
	m.chapters := make([]Chapter)
	tchap *Chapter
	doc,_ := goquery.NewDocument(fmt.Sprintf("%s%s",m.parent.base_path,m.rel_path))
	doc.Find("table#listing a").Each(func(i int, s *goquery.Selection){
		link,_ = s.Attr("href")
		tchap = make(Chapter)
		tchap.number = i
		tchap.starting_path = link
		tchap.parent = m.parent
		m.chapters[i] = tchap
	})
}

func (m *Manga) download_chapter_list(){
	// Calls chapter_list_$, with the links to the first page of each chapter,
	// calls Chapter.download_chapter(#) on each string of the returned slice
	// For now this only works with mangareader
	for _,chap := range m.chapters {
		go chap.download_chapter()
	}
	for _,chap := range m.chapters {
		<- chap.done
	}
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
	m.parent = c
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
