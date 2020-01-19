// 入口文件
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/xmluozp/creinox_server/driver"
	"github.com/xmluozp/creinox_server/routes"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

const PORT = "8000"

// 创建一个空的slice叫roles. Role类型
// var roles []models.Role
var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {

	db = driver.ConnectDB()
	// controller := controllers.Controller{}

	// 建立一个router
	router := mux.NewRouter()

	// 监听 {这里面是param}

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	// router.Handle("/", http.FileServer(http.Dir("./static")))

	// router.Use(static.Serve("/", static.LocalFile("./views", true)))

	routes.Routing(router)
	// router.HandleFunc("/api/role", controller.GetRoles(db)).Methods("GET") // 加个api避免混淆
	// router.HandleFunc("/api/role/{id}", controller.GetRole(db)).Methods("GET")
	// router.HandleFunc("/api/role", controller.AddRole(db)).Methods("POST")
	// router.HandleFunc("/api/role", controller.UpdateRole(db)).Methods("PUT")
	// router.HandleFunc("/api/role/{id}", controller.DeleteRole(db)).Methods("DELETE")

	fmt.Println("Server is running at port ", PORT)
	//第一个是端口，第二个是响应端口用的function。这里是router
	log.Fatal(http.ListenAndServe(":8000", router))
}
