package driver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xmluozp/creinox_server/initial"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {

	// cfg, err := goconfig.LoadConfigFile("conf.ini")

	// if err != nil {
	// 	panic("错误，找不到conf.ini配置文件")
	// }

	address, username, password, database := initial.GetMySql()

	connectString := fmt.Sprintf("%s:%s@%s/%s?parseTime=true&loc=Local", username, password, address, database)

	// db, err := sql.Open("mysql", "creinox:123456@/creinox?parseTime=true&loc=Local")
	db, err := sql.Open("mysql", connectString)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	// db.SetConnMaxLifetime(time.Second * 10)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}
