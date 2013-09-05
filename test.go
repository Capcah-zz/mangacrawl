package main

import(
	"fmt"
	"strconv"
	"archive/tar"
	"compress/gzip"
	"os"
	"io"
	"strings"
	"net/http"
	"bytes"
)

type Chapter struct {
	number int
	pages int
	title string
	tar_writer *tar.Writer 
}

func (chap *Chapter) download_page(path string, ret_chan *chan io.ReadCloser) {
	content,_ := http.Get(path)
	*ret_chan <- content.Body
	fmt.Println(path)
}

func (chap *Chapter) store_chapter(pages_data []io.ReadCloser) {
	var header *tar.Header
	buffer := new(bytes.Buffer)
	var written_size int64
	for i,data:= range pages_data {
		written_size,_ = io.Copy(buffer,data)
		header = &tar.Header{
			Name: fmt.Sprintf("./%d",i),
			Size: written_size,
			Mode: 0644,
		}
		chap.tar_writer.WriteHeader(header)
		io.Copy(chap.tar_writer,buffer)
	}
}

func (chap *Chapter) download_chapter(chapter_path string)  {
	//chapter path is made of i{rand 1-999}.mangareader.net/name/chapter/name-weirdnum.jpg
	//this function should generate the random {1-999} and for each page in the chapter,
	//add 2 to weirdnum and download it
	fd,_ := os.Create(chap.title)
	defer fd.Close()
	gw := gzip.NewWriter(fd)
	defer gw.Close()
	chap.tar_writer = tar.NewWriter(gw)
	defer chap.tar_writer.Close()
	sweirdnum,eweirdnum := strings.Index(chapter_path,"-"),strings.LastIndex(
									chapter_path,".")
	weirdnum,_ := strconv.Atoi(chapter_path[sweirdnum+1:eweirdnum])
	base_dl_path := chapter_path[0:sweirdnum+1]
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
}

func main(){
	// blob,_ := http.Get("http://i32.mangareader.net/claymore/142/claymore-4422419.jpg")
	// contents, _ := ioutil.ReadAll(blob.Body)
	// ioutil.WriteFile("image.jpg",contents,0664)
	// fmt.Println("Here!\n")
	chap := &Chapter{14,10,"testfile.tar.gz",nil}
	chap.download_chapter("http://i998.mangareader.net/claymore/142/claymore-4422419.jpg")
}
