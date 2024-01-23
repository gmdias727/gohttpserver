package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
)

func DatabaseInit() {
	log.Println("Deleting all files ending with .db")
	os.Remove("./database/*.db")

	log.Println("Creating a books.db file at database directory")
	file, err := os.Create("./database/books.db")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	createBooksTable()
	showAllBooks()
}

func createBooksTable() {
	db, err := sql.Open("sqlite3", "./database/books.db")

	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	stmt := `
	DROP TABLE IF EXISTS books;
	CREATE TABLE books (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		isbn TEXT NOT NULL UNIQUE
	);
	INSERT INTO books(title, isbn) VALUES('Getulio, O Cavaleiro', '123-456-789');
	INSERT INTO books(title, isbn) VALUES('Giraldo, O Poeta', '789-456-123');
	`

	log.Println("Executing the create books table query")
	_, err = db.Exec(stmt)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Books table created successfully")
}

func showAllBooks() {
	db, err := sql.Open("sqlite3", "./database/books.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err.Error())
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID \tName \tISBN")

	defer rows.Close()

	for rows.Next() {
		var (
			id    int
			title string
			isbn  string
		)

		err = rows.Scan(&id, &title, &isbn)

		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", id, title, isbn)
	}
	w.Flush()
}
