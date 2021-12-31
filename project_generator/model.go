package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"
)

type name struct {
	Name string `json:"name"`
}

func (n *name) getWords(DB *sql.DB) ([]name, error) {

	var err error
	var aNum int
	aRows := DB.QueryRow("SELECT COUNT (DISTINCT adjective) FROM words WHERE adjective IS NOT NULL")
	err = aRows.Scan(&aNum)
	if err != nil {
		log.Fatal(err)
	}

	var nNum int
	nRows := DB.QueryRow("SELECT COUNT (DISTINCT noun) FROM words WHERE noun IS NOT NULL")
	err = nRows.Scan(&nNum)
	if err != nil {
		log.Fatal(err)
	}

	min := 2
	rand.Seed(time.Now().UnixNano())

	var adjective string
	aRand := rand.Intn(aNum-min+1) + min
	aWord := DB.QueryRow("SELECT adjective FROM words where id=$1", aRand)
	err = aWord.Scan(&adjective)
	if err != nil {
		log.Fatal(err)
	}

	var noun string
	nRand := rand.Intn(nNum-min+1) + min
	nWord := DB.QueryRow("SELECT noun FROM words where id=$1", nRand)
	err = nWord.Scan(&noun)
	if err != nil {
		log.Fatal(err)
	}

	names := []name{}
	var projectName string = adjective + "-" + noun

	fcheck, err := os.ReadFile("name")
	if err != nil {
		log.Fatal(err)
	}

	result := string(fcheck) == projectName

	if result {
		return names, err
	} else {
		file, err := os.Create("name")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(projectName)
		if err != nil {
			log.Fatal(err)
		}
		var n name
		n.Name = projectName
		names = append(names, n)
		return names, nil
	}
}
