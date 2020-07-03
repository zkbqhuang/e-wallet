package config

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var err error

// Connect init
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/e-wallet")
	if err != nil {
		return nil, err
	}
	fmt.Println("Connection success")

	return db, nil
}
