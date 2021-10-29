package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Name struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (n *Name) getName(db *sql.DB) {

	var aNum int
	aRows := db.QueryRow("SELECT COUNT (DISTINCT adjective) FROM words")
	aRows.Scan(aNum)
	defer db.Close()

	var nNum int
	nRows := db.QueryRow("SELECT COUNT (DISTINCT nouns) FROM words")
	nRows.Scan(&nNum)
	defer db.Close()

	min := 1
	rand.Seed(time.Now().UnixNano())

	var adjective string
	aRand := rand.Intn(aNum+min) + min
	aWord := db.QueryRow("SELECT $1 FROM noun", aRand)
	aWord.Scan(adjective)
	defer db.Close()

	var noun string
	nRand := rand.Intn(nNum+min) + min
	nWord := db.QueryRow("SELECT $1 FROM adjective", nRand)
	nWord.Scan(noun)
	defer db.Close()

	file, fileErr := os.Create("name")
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
	fmt.Fprintf(file, "%v\n", projectName)
	var projectName string = adjective + "-" + noun

}
