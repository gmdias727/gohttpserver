package main

import (
	"fmt"

	database "github.com/gohttpserver/database"
)

func main() {
	fmt.Println("Main function starting...")
	database.DatabaseInit()
}
