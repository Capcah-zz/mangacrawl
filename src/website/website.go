package website

import(
	"fmt"
	"net/http"
)

func main(){
	server = &http.Server {
		Addr:		":3000",
		Hander:	
	}
	http.Handle("/files/",http.FileServer(http.Dir("$GOPATH/")))
}
