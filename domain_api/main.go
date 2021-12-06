package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8010")
	fmt.Println("in there")
}
