package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	bankaccountController "github.com/xmluozp/creinox_server/controllers/bankAccount"
	categoryController "github.com/xmluozp/creinox_server/controllers/category"
	commonitemController "github.com/xmluozp/creinox_server/controllers/commonItem"
	companyController "github.com/xmluozp/creinox_server/controllers/company"
	productController "github.com/xmluozp/creinox_server/controllers/product"

	testController "github.com/xmluozp/creinox_server/controllers/test"

	imageController "github.com/xmluozp/creinox_server/controllers/imagedata"
	regionController "github.com/xmluozp/creinox_server/controllers/region"
	roleController "github.com/xmluozp/creinox_server/controllers/role"
	rostercontactController "github.com/xmluozp/creinox_server/controllers/rosterContact"
	userController "github.com/xmluozp/creinox_server/controllers/user"
)

func Routing(router *mux.Router, db *sql.DB) {

	// ------------ test
	testController := testController.Controller{}
	router.HandleFunc("/api/test/{v}", testController.Test(db)).Methods("GET") // 加个api避免混淆

	// ------------ role
	roleController := roleController.Controller{}
	router.HandleFunc("/api/role", roleController.GetItems(db)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/role/{id}", roleController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/role", roleController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/role", roleController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/role/{id}", roleController.DeleteItem(db)).Methods("DELETE")

	// ------------ user
	userController := userController.Controller{}
	router.HandleFunc("/api/user", userController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/user/{id}", userController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/user", userController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/user", userController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/user/{id}", userController.DeleteItem(db)).Methods("DELETE")
	router.HandleFunc("/api/user/login", userController.Login(db)).Methods("POST")

	// ------------ commonitem
	commonitemController := commonitemController.Controller{}
	router.HandleFunc("/api/commonitem", commonitemController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/commonitem/{id}", commonitemController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/commonitem", commonitemController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/commonitem", commonitemController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/commonitem/{id}", commonitemController.DeleteItem(db)).Methods("DELETE")

	router.HandleFunc("/api/commonitem_dropDown", commonitemController.GetItems_DropDown(db)).Methods("GET")

	// ------------ image
	imageController := imageController.Controller{}
	router.HandleFunc("/api/image", imageController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/image/{id}", imageController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/image", imageController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/image/{folder_id}", imageController.AddItems(db)).Methods("POST")
	router.HandleFunc("/uploads/{path}", imageController.Show(db)).Methods("GET")

	router.HandleFunc("/api/image", imageController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/image_delete", imageController.DeleteItems(db)).Methods("PUT")

	// ------------ company
	companyController := companyController.Controller{}
	router.HandleFunc("/api/company", companyController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/company/{id}", companyController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/company", companyController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/company", companyController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/company/{id}", companyController.DeleteItem(db)).Methods("DELETE")

	// ------------ company: rostercontactController
	rostercontactController := rostercontactController.Controller{}
	router.HandleFunc("/api/rostercontact", rostercontactController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/rostercontact/{id}", rostercontactController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/rostercontact", rostercontactController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/rostercontact", rostercontactController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/rostercontact/{id}", rostercontactController.DeleteItem(db)).Methods("DELETE")

	// ------------ company: bankAccount
	bankaccountController := bankaccountController.Controller{}
	router.HandleFunc("/api/bankaccount", bankaccountController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/bankaccount/{id}", bankaccountController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/bankaccount", bankaccountController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/bankaccount", bankaccountController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/bankaccount/{id}", bankaccountController.DeleteItem(db)).Methods("DELETE")

	// ------------ region
	regionController := regionController.Controller{}
	router.HandleFunc("/api/region", regionController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/region/{id}", regionController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/region", regionController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/region", regionController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/region/{id}", regionController.DeleteItem(db)).Methods("DELETE")

	// ------------ category
	categoryController := categoryController.Controller{}
	router.HandleFunc("/api/category", categoryController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/category/{id}", categoryController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/category", categoryController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/category", categoryController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/category/{id}", categoryController.DeleteItem(db)).Methods("DELETE")

	// ------------ product
	productController := productController.Controller{}
	router.HandleFunc("/api/product", productController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/product/{id}", productController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/product", productController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/product", productController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/product/{id}", productController.DeleteItem(db)).Methods("DELETE")

	router.HandleFunc("/api/product_dropDown", productController.GetItems_DropDown(db)).Methods("GET")
}
