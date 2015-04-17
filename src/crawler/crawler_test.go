package crawler

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func find_string(s []string, t string) int {
	for i:=0; i< len(s); i++{
		if (s[i] == t){
			return i
		}
	}
	return -1
}

func TestCrawler_manga_list_create_entry(t *testing.T){
	c := &Crawler{
		last_crawl: time.Now(),
		website: "mangareader",
		nmangas: 0,
		base_path: "http://mangareader.net",
	}
	mlist := c.manga_list();
	fmt.Println(mlist)
	pindex := find_string(mlist,"/shingeki-no-kyojin")
	if (pindex == -1){
		t.Error("could not find SnK entry on manga_list")
	}
	m := c.create_manga_entry(mlist[pindex])
	if (	m.alt_names 	!= "" 	||
		m.release_year 	!= 00 	||
		m.ongoing	!= true	||
		m.author	!= ""	||
		m.artist	!= ""	||
		m.normal	!= true	){
		t.Error("There is a wrong information")
	}
}

func TestManga_download_chapter_list(t *testing.T) {
	c := &Crawler{
		last_crawl: time.Now(),
		website: "mangareader",
		nmangas: 10,
		base_path: "http://mangareader.net",
	}
	mlist := c.manga_list()
	pindex := find_string(mlist,"/shingeki-no-kyojin")
	stub_manga := &Manga{
		rel_path: mlist[pindex],
		chapters: nil,
		parent: &Crawler{
			base_path: "http://mangareader.net",
		},
	}
	stub_manga.parent.initialize_semaphors(50)
	stub_manga.chapter_list_mangareader()
	stub_manga.download_chapter_list()
	exec.Command("rm", "'Shingeki no Kyojin'*")
}

func TestChapter(t *testing.T) {
	stub_chapter := &Chapter{
		number: 1,
		pages:  10,
		title:  "test_chapter",
		done:   make(chan bool, 2),
		parent: &Crawler{},
	}
	stub_chapter.parent.initialize_semaphors(100)
	stub_chapter.download_chapter("http://i996.mangareader.net/freezing/130/freezing-4436531.jpg")
	outp, _ := exec.Command("du", "test_chapter.tar").Output()
	if string(outp) != "2832	test_chapter.tar\n" {
		t.Error("download_chapter did not work as expected")
	}
	exec.Command("rm", "test_chapter.tar").Run()
}


