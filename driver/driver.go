package driver

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {

	db, err := sql.Open("mysql", "creinox:123456@creinox")
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}
