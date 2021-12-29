package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	a.Run(":8080")
	fmt.Println("in there")
}
