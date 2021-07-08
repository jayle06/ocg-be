package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	user     = "root"
	password = "12345678"
	host     = "127.0.0.1:3305"
	database = "adidas"
)

func Connect() (db *sql.DB) {
	dsn := user + ":" + password + "@tcp(" + host + ")/" + database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("Could not connect to database")
	}
	return db
}
