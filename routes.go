package routes

import (
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/controllers"
)

func Routing(router *mux.Router) {
	controller := controllers.Controller{}
	router.HandleFunc("/api/role", controller.GetRoles(db)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/role/{id}", controller.GetRole(db)).Methods("GET")
	router.HandleFunc("/api/role", controller.AddRole(db)).Methods("POST")
	router.HandleFunc("/api/role", controller.UpdateRole(db)).Methods("PUT")
	router.HandleFunc("/api/role/{id}", controller.DeleteRole(db)).Methods("DELETE")
}
