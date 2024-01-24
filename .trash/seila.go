//go:build exclude

package main

import (
	"io"
	"log"
	"net/http"

	db "github.com/gohttpserver/db"
)

type Book2 struct {
	ID        int32  `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
	ISBN      string `json:"isbn"`
}

func main2897() {
	db.Database()

	helloWorldHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/", helloWorldHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
