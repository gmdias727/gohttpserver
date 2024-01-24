package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/gohttpserver/database"
	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()
	router.GET("/", index)
	// router.POST("/newBook", handleNewBook)

	database.DatabaseInit()

	log.Fatal(http.ListenAndServe(":8000", router))
}

type Book struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Payload struct {
	Data interface{} `json:"data"`
}

// GET all books
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	var book []Book

	for query.Next() {
		var (
			id    int
			title string
		)

		err = query.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}

		book = append(book, Book{
			Id:    id,
			Title: title,
		})
	}

	if len(book) <= 0 {
		log.Fatal("Your book list is empty")
	}

	resp := &Payload{Data: &book}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Data)
}

func getAllBooks() []Book {
	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	var books []Book

	for query.Next() {
		var (
			id    int
			title string
		)

		err = query.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}

		books = append(books, Book{
			Id:    id,
			Title: title,
		})
	}
	fmt.Println(books)
	return books
}

// func createOneBook(title string) Book {
// 	db, err := sql.Open("sqlite3", "./database/books.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	stmt := `
// 	INSERT INTO books (title)
// 	VALUES('Giraldo, O Cavaleiro Parte 2')
// 	`
// 	log.Println("Querying a new book into books")
// 	_, err = db.Exec(stmt)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Query successfull")
// }

// TODO: Allow dynamic data written by user
// TODO: Save user data to database
func handleNewBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Make a Go slice into a valid JSON to be sent via http req
	book := map[string]string{"title": "John Doe"} // remove this hardcoded
	json_data, err := json.Marshal(book)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	log.Println(res)
}
