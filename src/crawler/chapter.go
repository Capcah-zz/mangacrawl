package crawler

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Chapter struct {
	number        int
	pages         int
	title         string
	starting_path string
	done          chan bool
	tar_writer    *tar.Writer
	parent        *Crawler
}

func (chap *Chapter) download_page(path string, ret_chan *chan io.ReadCloser) {
	<-chap.parent.semaphores
	content := http_download_safe(path)
	*ret_chan <- content.Body
	chap.parent.semaphores <- true
}

func (chap *Chapter) store_chapter(pages_data []io.ReadCloser) {
	var header *tar.Header
	buffer := new(bytes.Buffer)
	var written_size int64
	for i, data := range pages_data {
		written_size, _ = io.Copy(buffer, data)
		header = &tar.Header{
			Name:    fmt.Sprintf("./%s/%d.jpg", chap.title, i),
			Size:    written_size,
			ModTime: time.Now(),
			Mode:    0644,
		}
		chap.tar_writer.WriteHeader(header)
		io.Copy(chap.tar_writer, buffer)
	}
	chap.done <- true
}
func (chap *Chapter) download_chapter(chapter_path string) {
	//chapter path is made of i{rand 1-999}.mangareader.net/name/chapter/name-weirdnum.jpg
	//this function should generate the random {1-999} and for each page in the chapter,
	//add 2 to weirdnum and download it
	fmt.Println(chapter_path)
	fd, _ := os.Create(fmt.Sprintf("%s.tar", chap.title))
	defer fd.Close()
	chap.tar_writer = tar.NewWriter(fd)
	defer chap.tar_writer.Close()
	sweirdnum, eweirdnum := strings.LastIndex(chapter_path, "-"), strings.LastIndex(
		chapter_path, ".")
	weirdnum, _ := strconv.Atoi(chapter_path[sweirdnum+1 : eweirdnum])
	base_dl_path := chapter_path[0 : sweirdnum+1]
	pages_chanel := make([]chan io.ReadCloser, chap.pages)
	pages_data := make([]io.ReadCloser, chap.pages)
	for i, _ := range pages_chanel {
		pages_chanel[i] = make(chan io.ReadCloser)
	}
	for i, _ := range pages_data {
		go chap.download_page(fmt.Sprintf("%s%d.jpg", base_dl_path, weirdnum+2*i),
			&pages_chanel[i])
	}
	for i, _ := range pages_data {
		pages_data[i] = <-(pages_chanel[i])
	}
	chap.store_chapter(pages_data)
	fmt.Println("Finished chapter")
}
