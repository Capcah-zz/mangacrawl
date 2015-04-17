package crawler

import (
	"fmt"
	"time"
)

func main() {
	crawl := Crawler{time.Now(), "mangareader", 0, "http://mangareader.net", make(chan bool, 1000), make(map[string]Manga)}
	crawl.initialize_semaphors(1000)
	crawl.update_mangas()
	for crawl.nmangas > 0 {
		time.Sleep(time.Second)
	}
	for _, manga := range crawl.mangas {
		fmt.Printf("%s\n", manga.name)
	}
}
