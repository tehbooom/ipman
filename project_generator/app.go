package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	var err error
	response, _ := json.MarshalIndent(payload, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) getName(w http.ResponseWriter, r *http.Request) {

	var err error
	n := name{}
	names, errGet := n.getWords(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if errGet != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, names)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/project", a.getName).Methods("GET")
}
