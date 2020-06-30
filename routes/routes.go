package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	bankaccountController "github.com/xmluozp/creinox_server/controllers/bankAccount"
	categoryController "github.com/xmluozp/creinox_server/controllers/category"
	commodityController "github.com/xmluozp/creinox_server/controllers/commodity"
	commonitemController "github.com/xmluozp/creinox_server/controllers/commonItem"
	companyController "github.com/xmluozp/creinox_server/controllers/company"
	productController "github.com/xmluozp/creinox_server/controllers/product"
	productPurchaseController "github.com/xmluozp/creinox_server/controllers/productPurchase"

	sellContractController "github.com/xmluozp/creinox_server/controllers/sellContract"
	sellSubitemController "github.com/xmluozp/creinox_server/controllers/sellSubitem"

	buyContractController "github.com/xmluozp/creinox_server/controllers/buyContract"
	buySubitemController "github.com/xmluozp/creinox_server/controllers/buySubitem"

	mouldContractController "github.com/xmluozp/creinox_server/controllers/mouldContract"

	testController "github.com/xmluozp/creinox_server/controllers/test"
	textTemplateController "github.com/xmluozp/creinox_server/controllers/textTemplate"

	imageController "github.com/xmluozp/creinox_server/controllers/imagedata"
	portController "github.com/xmluozp/creinox_server/controllers/port"
	regionController "github.com/xmluozp/creinox_server/controllers/region"

	roleController "github.com/xmluozp/creinox_server/controllers/role"
	rostercontactController "github.com/xmluozp/creinox_server/controllers/rosterContact"
	userController "github.com/xmluozp/creinox_server/controllers/user"

	printController "github.com/xmluozp/creinox_server/controllers/printdata"
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

	router.HandleFunc("/api/userList", userController.GetItemsForLogin(db)).Methods("GET")

	// ------------ text template
	textTemplateController := textTemplateController.Controller{}
	router.HandleFunc("/api/texttemplate", textTemplateController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/texttemplate/{id}", textTemplateController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/texttemplate", textTemplateController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/texttemplate", textTemplateController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/texttemplate/{id}", textTemplateController.DeleteItem(db)).Methods("DELETE")

	router.HandleFunc("/api/texttemplate_all", textTemplateController.GetItems_Template(db)).Methods("GET")

	// ------------ print
	printController := printController.Controller{}
	router.HandleFunc("/api/printFolder/{templateFolder}", printController.GetItems(db)).Methods("GET")

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

	// ------------ port
	portController := portController.Controller{}
	router.HandleFunc("/api/port", portController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/port/{id}", portController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/port", portController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/port", portController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/port/{id}", portController.DeleteItem(db)).Methods("DELETE")

	// ------------ category
	categoryController := categoryController.Controller{}
	router.HandleFunc("/api/category", categoryController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/category/{id}", categoryController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/category", categoryController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/category", categoryController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/category/{id}", categoryController.DeleteItem(db)).Methods("DELETE")

	// ------------ sell contract
	sellContractController := sellContractController.Controller{}
	router.HandleFunc("/api/sellcontract", sellContractController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/sellcontract/{id}", sellContractController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/sellcontract", sellContractController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/sellcontract", sellContractController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/sellcontract/{id}", sellContractController.DeleteItem(db)).Methods("DELETE")
	router.HandleFunc("/api/sellcontract_print/{id}/{templateFolder}/{template}/{printFormat}", sellContractController.Print(db)).Methods("GET")

	// customized
	router.HandleFunc("/api/sellcontract_getlast", sellContractController.GetLast(db)).Methods("GET")

	// ------------ sell subitem
	sellSubitemController := sellSubitemController.Controller{}
	router.HandleFunc("/api/sellsubitem", sellSubitemController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/sellsubitem/{id}", sellSubitemController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/sellsubitem", sellSubitemController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/sellsubitem", sellSubitemController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/sellsubitem/{id}", sellSubitemController.DeleteItem(db)).Methods("DELETE")

	// ------------ buy contract
	buyContractController := buyContractController.Controller{}
	router.HandleFunc("/api/buycontract", buyContractController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/buycontract/{id}", buyContractController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/buycontract", buyContractController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/buycontract", buyContractController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/buycontract/{id}", buyContractController.DeleteItem(db)).Methods("DELETE")
	router.HandleFunc("/api/buycontract_print/{id}/{templateFolder}/{template}/{printFormat}", buyContractController.Print(db)).Methods("GET")

	// customized
	router.HandleFunc("/api/buycontract_getlast", buyContractController.GetLast(db)).Methods("GET")

	// ------------ buy subitem
	buySubitemController := buySubitemController.Controller{}
	router.HandleFunc("/api/buysubitem", buySubitemController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/buysubitem/{id}", buySubitemController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/buysubitem", buySubitemController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/buysubitem", buySubitemController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/buysubitem/{id}", buySubitemController.DeleteItem(db)).Methods("DELETE")

	// ------------ mould contract
	mouldContractController := mouldContractController.Controller{}
	router.HandleFunc("/api/mouldcontract", mouldContractController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/mouldcontract/{id}", mouldContractController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/mouldcontract", mouldContractController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/mouldcontract", mouldContractController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/mouldcontract/{id}", mouldContractController.DeleteItem(db)).Methods("DELETE")
	router.HandleFunc("/api/mouldcontract_print/{id}/{templateFolder}/{template}/{printFormat}", mouldContractController.Print(db)).Methods("GET")

	// customized
	router.HandleFunc("/api/mouldcontract_getlast", mouldContractController.GetLast(db)).Methods("GET")

	// ------------ product
	productController := productController.Controller{}
	router.HandleFunc("/api/product", productController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/product/{id}", productController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/product", productController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/product", productController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/product/{id}", productController.DeleteItem(db)).Methods("DELETE")

	router.HandleFunc("/api/product_dropDown", productController.GetItems_DropDown(db)).Methods("GET")
	router.HandleFunc("/api/product_dropDown_sellContract", productController.GetItems_DropDown_sellContract(db)).Methods("GET")
	router.HandleFunc("/api/product_dropDown_sellSubitem", productController.GetItems_DropDown_sellSubitem(db)).Methods("GET")

	router.HandleFunc("/api/product_component", productController.GetComponents(db)).Methods("GET")
	router.HandleFunc("/api/product_component/{parent_id}/{child_id}", productController.Assemble(db)).Methods("POST")
	router.HandleFunc("/api/product_component/{parent_id}/{child_id}", productController.Disassemble(db)).Methods("DELETE")

	// commodity_product
	router.HandleFunc("/api/commodity_getproduct", productController.GetItems_ByCommodity(db)).Methods("GET")

	// ------------ product purchase
	productPurchaseController := productPurchaseController.Controller{}
	router.HandleFunc("/api/productPurchase", productPurchaseController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/productPurchase/{id}", productPurchaseController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/productPurchase", productPurchaseController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/productPurchase", productPurchaseController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/productPurchase/{id}", productPurchaseController.DeleteItem(db)).Methods("DELETE")

	router.HandleFunc("/api/productPurchase_companySearch", productPurchaseController.GetItems_GroupByCompany(db)).Methods("GET")
	router.HandleFunc("/api/productPurchase_historySearch", productPurchaseController.GetItems_History(db)).Methods("GET")
	router.HandleFunc("/api/productPurchase_byProductId/{id}", productPurchaseController.GetItem_ByProductId(db)).Methods("GET")

	// ------------ commodity
	commodityController := commodityController.Controller{}
	router.HandleFunc("/api/commodity", commodityController.GetItems(db)).Methods("GET")
	router.HandleFunc("/api/commodity/{commodity_id}", commodityController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/commodity/{commodity_id}/{product_id}", commodityController.GetItem(db)).Methods("GET")
	router.HandleFunc("/api/commodity", commodityController.AddItem(db)).Methods("POST")
	router.HandleFunc("/api/commodity", commodityController.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/api/commodity/{id}", commodityController.DeleteItem(db)).Methods("DELETE")

	router.HandleFunc("/api/commodity_byproduct", commodityController.GetItems_ByProduct(db)).Methods("GET")
	// commodity_product
	router.HandleFunc("/api/commodity_byproduct/{commodity_id}/{product_id}", commodityController.Assemble(db)).Methods("POST")
	router.HandleFunc("/api/commodity_byproduct/{commodity_id}/{product_id}", commodityController.Disassemble(db)).Methods("DELETE")

}
