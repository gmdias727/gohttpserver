package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Should I make a test for this code? I think yes...
// Should display 3.44.0
func DisplaySqliteVersion() {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sqlite3 version:", version)
}
