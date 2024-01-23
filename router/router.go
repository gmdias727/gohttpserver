//go:build exclude

package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// bookstore["123"] = &Book{
// 	Title: "Silence of the Lambs",
// 	Isbn:  "123",
// }

func Router() {
	port := ":8000"
	router := httprouter.New()
	router.GET("/", IndexHandler)
	// router.GET("/hello/:book", BookIndexHandler)

	log.Printf("Starting server at http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

type Book struct {
	Title string `json:"title"`
	Isbn  string `json:"isbn"`
}

type JsonResponse struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

var bookstore = make(map[string]*Book)

func IndexHandler(w http.ResponseWriter, resp *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// GET all books
// func BookIndexHandler(w http.ResponseWriter, res *http.Request, ps httprouter.Params) {
// 	db, err := sql.Open("sqlite3", "books.db")

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer db.Close()

// 	// for _, book := range bookstore { // loop through all existing books
// 	// 	books = append(books, book)
// 	// }

// 	response := &JsonResponse{Data: &books}                           // Add the retrieved books to Json
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8") // Set this request as "application/json"
// 	w.WriteHeader(http.StatusOK)                                      // HTTP Status = 200 OK

// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		panic(err)
// 	}
// }
