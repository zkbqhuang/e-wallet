package main

import (
	"database/sql"
	"privyid/config"
	"privyid/service/delivery"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func main() {
	// Init connection ke DB
	config.Connect()

	// Init router
	delivery.Router()

}
