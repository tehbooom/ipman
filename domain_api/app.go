package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) getDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uri := vars["domain"]

	d := domain{Domain: uri}
	if err := d.getDomain(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Domain not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) getOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	op := vars["operation"]

	d := domain{Operation: op}
	domains, err := d.getOperation(a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Operation not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, domains)
}

func (a *App) getDomains(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count, _ := strconv.Atoi(vars["count"])

	var errInt error
	if count <= 0 {
		err := fmt.Sprintf("Invalid count: %d", count)
		errInt = errors.New(err + ". Count must be above zero.")
	}

	domains, errGet := getDomains(a.DB, count)
	if errInt != nil {
		respondWithError(w, http.StatusInternalServerError, errInt.Error())
		return
	} else if errGet != nil {
		respondWithError(w, http.StatusInternalServerError, errGet.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, domains)
}

func (a *App) createDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uri := vars["domain"]

	d := domain{Domain: uri}

	if err := d.createDomain(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, d)
}

func (a *App) updateDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid domain")
		return
	}

	var d domain
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	d.ID = id

	if err := d.updateDomain(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) deleteDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Domain ID")
		return
	}

	d := domain{ID: id}
	if err := d.deleteDomain(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/getDomains/{count}", a.getDomains).Methods("GET")
	a.Router.HandleFunc("/create/{domain}", a.createDomain).Methods("POST")
	a.Router.HandleFunc("/status/domain/{domain}", a.getDomain).Methods("GET")
	a.Router.HandleFunc("/status/operation/{operation}", a.getOperation).Methods("GET")
	a.Router.HandleFunc("/update/{domain}", a.updateDomain).Methods("PUT")
	a.Router.HandleFunc("/delete/{domain}", a.deleteDomain).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
