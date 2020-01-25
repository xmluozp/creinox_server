// 入口文件
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/driver"
	"github.com/xmluozp/creinox_server/routes"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

const PORT = "8000"

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {

	db = driver.ConnectDB()

	auth.JwtKey = []byte{1}
	// rand.Read(auth.JwtKey)

	// 建立一个router
	router := mux.NewRouter()

	// 监听 {这里面是param}

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	// router.Handle("/", http.FileServer(http.Dir("./static")))

	// router.Use(static.Serve("/", static.LocalFile("./views", true)))

	routes.Routing(router, db)

	fmt.Println("Server is running at port ", PORT)
	//第一个是端口，第二个是响应端口用的function。这里是router
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
