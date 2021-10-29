package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/subnets", AllSubnets).Methods("GET")
	myRouter.HandleFunc("/subnets/allocations/{subnet}/{ips}/{operation}", NewIP).Methods("POST")
	myRouter.HandleFunc("/subnets/allocations/{operation}", DeleteIP).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	InitialMigration()
	handleRequests()
}
