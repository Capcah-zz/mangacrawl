package crawler

import (
	"fmt"
	"os/exec"
	"testing"
)

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

func TestManga_download_chapter_list(t *testing.T) {
	stub_manga := &Manga{
		rel_path: "/shingeki-no-kyojin",
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
