package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	companyController "github.com/xmluozp/creinox_server/controllers/company"
	regionController "github.com/xmluozp/creinox_server/controllers/region"
	roleController "github.com/xmluozp/creinox_server/controllers/role"
	userController "github.com/xmluozp/creinox_server/controllers/user"
)

func Routing(router *mux.Router, db *sql.DB) {

	// ------------ role
	roleController := roleController.Controller{}
	router.HandleFunc("/api/role", roleController.GetItems(db)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/role/{id}", roleController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/role", roleController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/role", roleController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/role/{id}", roleController.DeleteItem(db)).Methods("DELETE")

	// ------------ user
	userController := userController.Controller{}
	router.HandleFunc("/api/user", userController.GetItems(db)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/user/{id}", userController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/user", userController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/user", userController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/user/{id}", userController.DeleteItem(db)).Methods("DELETE")
	router.HandleFunc("/api/user/login", userController.Login(db)).Methods("POST")

	// ------------ company
	companyController := companyController.Controller{}
	router.HandleFunc("/api/company", companyController.GetItems(db)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/company/{id}", companyController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/company", companyController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/company", companyController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/company/{id}", companyController.DeleteItem(db)).Methods("DELETE")

	// ------------ region
	regionController := regionController.Controller{}
	router.HandleFunc("/api/region", regionController.GetItems(db)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/region/{id}", regionController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/region", regionController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/region", regionController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/region/{id}", regionController.DeleteItem(db)).Methods("DELETE")
}
