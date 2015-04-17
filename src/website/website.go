package main

import(
	"net/http"
	"html/template"
	"io/ioutil"
	"regexp"
	"fmt"
	"crawler"
)

func asset_handler(w http.ResponseWriter, r *http.Request){
	data,err := ioutil.ReadFile(r.URL.Path[6:])
	if err != nil{
		http.NotFound(w,r)
		return
	}
	w.Write(data)
}

var manga_list_regex = regexp.MustCompile(`/[a-zA-z_-]+`)

func chapter_list(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Chapter list not implemented")
}

func single_chapter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "single_chapter not implemented")
}

func manga_list_handler(w http.ResponseWriter, r *http.Request){
	match_list := manga_list_regex.FindAllString(r.URL.Path[6:],-1)
	switch len(match_list) {
	case 1:
		chapter_list(w,r)
	case 2:
		single_chapter(w,r)
	}
}

func home_handler(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("index.html")
	//Retrieve the last updated mangas
	t.Execute
}

func crawler_daemon(){

}

func main(){
	http.HandleFunc("/asset/",asset_handler)
	http.HandleFunc("/manga/",manga_list_handler)
	http.HandleFunc("/",home_handler)
	http.ListenAndServe(":3000",nil)
}
