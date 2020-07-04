package config

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var err error

// Connect init untuk memulai koneksi ke database
func Connect() (*sql.DB, error) {
	// detail entitiy yang diperlukan, dengan mencantumkan detail dari database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/e-wallet")
	if err != nil {
		return nil, err
	}
	// status jika sukses melakukan koneksi ke database
	fmt.Println("Connection success")

	return db, nil
}
