package main

import(
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

func main(){
	blob,_ := goquery.NewDocument("http://www.mangareader.net/247/love-celeb.html")
	blob.Find("div#mangaproperties table td").Each(func(i int, s *goquery.Selection){
		fmt.Printf("Index %d found: %s\n",i,s.Text())
	})
	fmt.Println("Here!\n")
}
