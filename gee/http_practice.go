package gee

import (
	"fmt"
	"net/http"
)

type TEngine struct {
}

func (e *TEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "hello world!")
	case "/hello":
		fmt.Fprintf(w, "I am sb!")
	default:
		fmt.Fprintf(w, "404 not found! %s", r.URL)
	}
}

//func main() {
//	engine := new(TEngine)
//	http.ListenAndServe(":9999", engine)
//}

//func main() {
//	c := 0
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "URL.Path = %q", r.URL.Path)
//	})
//	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "count is %d", c)
//		c++
//	})
//	log.Fatal(http.ListenAndServe(":9999", nil))
//}
