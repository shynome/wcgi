package main

import (
	"io"
	"net/http"

	"github.com/shynome/wcgi"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "world")
	})
	wcgi.Serve(nil)
}
