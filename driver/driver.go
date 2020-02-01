package driver

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {

	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}
	address, err := cfg.GetValue("mysql", "address")
	username, err := cfg.GetValue("mysql", "username")
	password, err := cfg.GetValue("mysql", "password")
	database, err := cfg.GetValue("mysql", "database")

	connectString := fmt.Sprintf("%s:%s@%s/%s?parseTime=true&loc=Local", username, password, address, database)

	// db, err := sql.Open("mysql", "creinox:123456@/creinox?parseTime=true&loc=Local")
	db, err := sql.Open("mysql", connectString)

	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}
