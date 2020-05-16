// 入口文件
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Unknwon/goconfig"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/driver"
	"github.com/xmluozp/creinox_server/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {

	// fetch config file
	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}
	port, err := cfg.GetValue("site", "port")

	db = driver.ConnectDB()

	auth.JwtKey = []byte{1}
	// rand.Read(auth.JwtKey)

	// 建立一个router
	router := mux.NewRouter()

	// 监听 {这里面是param}

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	// router.Handle("/", http.FileServer(http.Dir("./static")))

	// router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// router是指针，探进去被修改
	routes.Routing(router, db)

	fmt.Println("Server is running at port ", port)
	//第一个是端口，第二个是响应端口用的function。这里是router

	// handler := cors.Default().Handler(router) //解决cors用

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"}) // domain

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))

	// log.Fatal(http.ListenAndServe(":"+port, router))

	// log.Fatal(http.ListenAndServe(":"+port, handler))
}
