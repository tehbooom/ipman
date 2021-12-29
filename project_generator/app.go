package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	dsn := fmt.Sprintf("host=localhost port=5432 sslmode=disable user=%s password=%s dbname=%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error connecting to db: %v\n", err)
	}

	a.initializeDB()

	file, err := os.Create("name") // create file for previous name
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getName(w http.ResponseWriter, r *http.Request) {

	// grab number of rows in column adjective
	var aNum int
	aRows := a.DB.QueryRow("SELECT COUNT (DISTINCT adjective) FROM words")
	aRows.Scan(aNum)
	defer a.DB.Close()

	// grab number of rows in column noun
	var nNum int
	nRows := a.DB.QueryRow("SELECT COUNT (DISTINCT noun) FROM words")
	nRows.Scan(&nNum)
	defer a.DB.Close()

	// setting vars for getting an almost random number
	min := 1
	rand.Seed(time.Now().UnixNano())

	// select random noun
	var adjective string
	aRand := rand.Intn(aNum+min) + min
	aWord := a.DB.QueryRow("SELECT $1 FROM noun", aRand)
	aWord.Scan(adjective)
	defer a.DB.Close()

	// select random adjective
	var noun string
	nRand := rand.Intn(nNum+min) + min
	nWord := a.DB.QueryRow("SELECT $1 FROM adjective", nRand)
	nWord.Scan(noun)
	defer a.DB.Close()

	var projectName string = adjective + "-" + noun

	check, err := os.ReadFile("name")
	if err != nil {
		log.Fatal(err)
	}

	result := string(check) == projectName

	if result == true {
		respondWithJSON(w, http.StatusConflict, "Previous project name please run again")
	} else {
		file, err := os.Create("name")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(projectName)
		respondWithJSON(w, http.StatusOK, projectName)
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/project", a.getName).Methods("GET")
}
