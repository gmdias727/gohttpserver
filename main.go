package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	database "github.com/gohttpserver/database"
	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()
	router.GET("/", handleGetSelectedOnly)
	router.POST("/book", handleNewBook)
	router.PUT("/book/:id", handleUpdateBook)
	router.DELETE("/book/:id", handleDeleteOneBook)

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

type Routerer interface {
	index(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if len(book) == 0 {
		book = make([]Book, 0)
	}

	resp := &Payload{Data: &book}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Data)
}

// TODO: Allow dynamic data written by user
// TODO: Save user data to database
func handleNewBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var newBook Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newBook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt := "INSERT INTO books (title) VALUES(?);"

	_, err = db.Exec(stmt, newBook.Title)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New book saved successfully")
}

// UPDATE books SET title = "test" WHERE id = 1;
func handleUpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var updatedBook Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	convertId, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	stmt := `UPDATE books SET title = ? WHERE id = ?;`
	_, err = db.Exec(stmt, updatedBook.Title, convertId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Book with id:%s updated successfully", id)
}

func handleDeleteOneBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	convertId, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	// DELETE FROM artists_backup WHERE artistid = 1;
	stmt := `DELETE FROM books WHERE id = ?`

	_, err = db.Exec(stmt, convertId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Book containing id: %s has been deleted.", id)
}

// http://localhost:8000/?id=hello&name=1,41,37&surname=123
// Need to delete id=1,2,3
// How to separate the
func handleGetSelectedOnly(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	values := strings.Split(queryValues.Get("id"), ",")

	var listIntIDs []int

	for _, id := range values {
		intId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Empty string not allowed here", http.StatusBadRequest)
			return
		}

		listIntIDs = append(listIntIDs, intId)

		// fmt.Printf("Id: %s \n", id)
	}

	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cases:
	// SELECT FROM books WHERE id IN (?)
	// SELECT FROM books WHERE id IN (?, ?)
	// SELECT FROM books WHERE id IN (?, ?, ?...+)
	// TODO: IMPLEMENT THE SELECT FOR DYNAMIC ID AMOUNT
	stmt := `SELECT FROM books WHERE id IN ?`
	_, err = db.Exec(stmt, listIntIDs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Selected books: %d", listIntIDs)
}

// func handleDeleteSelectedOnly(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	getUrlParams := r.URL.Query()
// 	values := strings.Split(getUrlParams.Get("id"), ",")
// 	log.Println(values)
// }
