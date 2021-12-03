package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const nurl = "https://greenopolis.com/list-of-nouns/"
const aurl = "https://greenopolis.com/adjectives-list/"

func (a *App) initializeDB() {

	// regex
	re := regexp.MustCompile(`<li>(.*?)</li>`)

	// create table
	const table = ` CREATE TABLE [IF NOT EXISTS] words (
		id serial PRIMARY KEY,
		noun text
		adjective text
	)`
	a.DB.Exec(table)

	// nouns

	const nrow = `INSERT INTO words (
	noun
	)
	VALUES $1
	)`

	nresp, err := http.Get(nurl) // get contents of noun webpage
	if err != nil {
		log.Fatal(err)
	}
	defer nresp.Body.Close()
	nhtml, err := ioutil.ReadAll(nresp.Body)
	if err != nil {
		log.Fatal(err)
	}

	nmatches := re.FindAllStringSubmatch(string(nhtml), -1) // select all words from regex and insert into table
	for _, w := range nmatches {
		a.DB.Exec(nrow, w[1])
	}
	defer a.DB.Close()

	//adjectives
	const arow = `INSERT INTO words (
		adjective
	)
	VALUES $1
	)`

	aresp, err := http.Get(aurl) // get contents of adjectives webpage
	if err != nil {
		log.Fatal(err)
	}
	defer aresp.Body.Close()
	ahtml, err := ioutil.ReadAll(aresp.Body)
	if err != nil {
		log.Fatal(err)
	}

	amatches := re.FindAllStringSubmatch(string(ahtml), -1) // select all words from regex and insert into table
	for _, w := range amatches {
		a.DB.Exec(arow, w[1])
	}
	defer a.DB.Close()
}
