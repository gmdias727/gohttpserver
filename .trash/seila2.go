//go:build exclude

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	// _ stands for blank identifier
	// The init function of the package is called
	routes "github.com/gohttpserver/Routes"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	Title string `json:"title"`
	Isbn  string `json:"isbn"`
}

func main() {
	routes.Router()
	// createBooksTable()

}

// Should query one book from database given id
// Isso aqui ta bem ruim eu acho
func queryOneBook(id int) {
	db, err := sql.Open("sqlite3", "./db/books.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM books WHERE id = ?")

	defer stmt.Close()

	// var id int
	var title string
	var isbn string
	cid := 2

	err = stmt.QueryRow(cid).Scan(&id, &title, &isbn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d %s %s\n", id, title, isbn)
}

func showAllBooksAsTable() {
	db, err := sql.Open("sqlite3", "./db/books.db")

	if err != nil {
		log.Fatal("Error while Open() attempt: ", err)
	}

	defer db.Close()

	query := `SELECT * FROM books`

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal("Error while Query() attempt", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID \tName \tISBN")

	defer rows.Close()

	for rows.Next() {
		var (
			id   int32
			name string
			isbn string
		)

		err = rows.Scan(&id, &name, &isbn)

		if err != nil {
			log.Fatal(err)
		}

		// resultList := []interface{}{id, name, isbn}

		// fmt.Printf("resultList: %s", resultList)
		fmt.Fprintf(w, "%d\t%s\t%s\n", id, name, isbn)
	}
	w.Flush()
}

func createBooksTable() {
	db, err := sql.Open("sqlite3", "./db/books.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	stmt := `
	DROP TABLE IF EXISTS books;
	CREATE TABLE books(id INTEGER PRIMARY KEY, name TEXT, isbn TEXT);
	INSERT INTO books(name, isbn) VALUES("blob", "123-456-789-0");
	`

	_, err = db.Exec(stmt)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`Table "Books" created successfully!\n`)
}

func customInsertion(name string, isbn string) {
	db, err := sql.Open("sqlite3", "./db/books.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	stmt := fmt.Sprintf(`INSERT INTO books(name, isbn) VALUES("%s","%s");`, name, isbn)

	_, err = db.Exec(stmt)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(`Name: %s and ISBN: %s has been added successfully to the database\n`, name, isbn)
}
