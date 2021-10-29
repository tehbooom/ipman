package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgtype"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// const specifying db connection
const (
	host     = "localhost"
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "initdb"
	port     = 5432
	sslmode  = "disable"
	TimeZone = "ETC"
)

// create new tpye subnet as a struct
type allocation struct {
	gorm.Model
	subnet    pgtype.CIDR
	ips       pgtype.Inet
	operation string
}

// sets database connection options as varaible
var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, TimeZone)

var db *gorm.DB

func InitialMigration() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&allocation{})
}

// grabs all subnets available in database
func AllSubnets(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var subnets []allocation
	db.Find(&subnets)
	json.NewEncoder(w).Encode(subnets)
}

// takes POST request of subent and number of ips requesting and updates the table with newlly allocated IPs
func NewIP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "NewIP endpoint Hit")

}

// based on opeartion name deletes all columns with value of operation name
func DeleteIP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DeleteIP endpoint Hit")
}
