package routes

import (
	"github.com/gorilla/mux"
	applicationController "github.com/xmluozp/creinox_server/controllers/application"
	bankaccountController "github.com/xmluozp/creinox_server/controllers/bankAccount"
	categoryController "github.com/xmluozp/creinox_server/controllers/category"
	commodityController "github.com/xmluozp/creinox_server/controllers/commodity"
	commonitemController "github.com/xmluozp/creinox_server/controllers/commonItem"
	companyController "github.com/xmluozp/creinox_server/controllers/company"
	"github.com/xmluozp/creinox_server/models"

	financialaccountController "github.com/xmluozp/creinox_server/controllers/financialAccount"
	financialledgerController "github.com/xmluozp/creinox_server/controllers/financialLedger"
	financialTransactionController "github.com/xmluozp/creinox_server/controllers/financialTransaction"
	financialVoucherController "github.com/xmluozp/creinox_server/controllers/financialVoucher"
	paymentRequestController "github.com/xmluozp/creinox_server/controllers/paymentRequest"
	productController "github.com/xmluozp/creinox_server/controllers/product"
	productPurchaseController "github.com/xmluozp/creinox_server/controllers/productPurchase"

	expressOrderController "github.com/xmluozp/creinox_server/controllers/expressOrder"

	orderformController "github.com/xmluozp/creinox_server/controllers/orderForm"
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
	userLogController "github.com/xmluozp/creinox_server/controllers/userLog"

	printController "github.com/xmluozp/creinox_server/controllers/printdata"
)

func Routing(router *mux.Router, mydb models.MyDb) {

	// ------------ test
	testController := testController.Controller{}
	router.HandleFunc("/api/test/{v}", testController.Test(mydb)).Methods("GET")                      // 加个api避免混淆
	router.HandleFunc("/api/testApp/{v}", testController.TestApp(mydb)).Methods("POST")               // 加个api避免混淆
	router.HandleFunc("/api/testAppReceive/{v}", testController.TestAppReceive(mydb)).Methods("POST") // 加个api避免混淆
	router.HandleFunc("/api/testTx/{v}", testController.TestTx(mydb.Db)).Methods("GET")               // 加个api避免混淆

	// ------------ application 申请(泛申请，暂未开发)
	applicationController := applicationController.Controller{}
	router.HandleFunc("/api/application", applicationController.GetItems(mydb)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/application/{id}", applicationController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/application", applicationController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/application", applicationController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/application/{id}", applicationController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ role
	roleController := roleController.Controller{}
	router.HandleFunc("/api/role", roleController.GetItems(mydb)).Methods("GET") // 加个api避免混淆
	router.HandleFunc("/api/role/{id}", roleController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/role", roleController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/role", roleController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/role/{id}", roleController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ user
	userController := userController.Controller{}
	router.HandleFunc("/api/user", userController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/user/{id}", userController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/user", userController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/user", userController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/user/{id}", userController.DeleteItem(mydb)).Methods("DELETE")
	router.HandleFunc("/api/user/login", userController.Login(mydb)).Methods("POST")

	router.HandleFunc("/api/userList", userController.GetItemsForLogin(mydb)).Methods("GET")

	// ------------ userLog 用户操作记录
	userLogController := userLogController.Controller{}
	router.HandleFunc("/api/userLog", userLogController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/userLog/{id}", userLogController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/userLog", userLogController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/userLog", userLogController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/userLog/{id}", userLogController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/userLog_delete", userLogController.DeleteItems(mydb)).Methods("PUT")

	// ------------ text template
	textTemplateController := textTemplateController.Controller{}
	router.HandleFunc("/api/texttemplate", textTemplateController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/texttemplate/{id}", textTemplateController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/texttemplate", textTemplateController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/texttemplate", textTemplateController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/texttemplate/{id}", textTemplateController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/texttemplate_all", textTemplateController.GetItems_Template(mydb)).Methods("GET")

	// ------------ print
	printController := printController.Controller{}
	router.HandleFunc("/api/printFolder/{templateFolder}", printController.GetItems(mydb)).Methods("GET")

	// ------------ commonitem
	commonitemController := commonitemController.Controller{}
	router.HandleFunc("/api/commonitem", commonitemController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/commonitem/{id}", commonitemController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/commonitem", commonitemController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/commonitem", commonitemController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/commonitem/{id}", commonitemController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/commonitem_dropDown", commonitemController.GetItems_DropDown(mydb)).Methods("GET")

	// ------------ image
	imageController := imageController.Controller{}
	router.HandleFunc("/api/image", imageController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/image/{id}", imageController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/image", imageController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/image/{folder_id}", imageController.AddItems(mydb)).Methods("POST")
	router.HandleFunc("/uploads/{path}", imageController.Show(mydb)).Methods("GET")

	router.HandleFunc("/api/image", imageController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/image_delete", imageController.DeleteItems(mydb)).Methods("PUT")

	// ------------ company
	companyController := companyController.Controller{}
	router.HandleFunc("/api/company", companyController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/company/{id}", companyController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/company", companyController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/company", companyController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/company/{id}", companyController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/companyGetCode/{companyType}/{keyWord}", companyController.GetRow_byCode(mydb)).Methods("GET")

	// ------------ company: rostercontactController
	rostercontactController := rostercontactController.Controller{}
	router.HandleFunc("/api/rostercontact", rostercontactController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/rostercontact/{id}", rostercontactController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/rostercontact", rostercontactController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/rostercontact", rostercontactController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/rostercontact/{id}", rostercontactController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ company: bankAccount
	bankaccountController := bankaccountController.Controller{}
	router.HandleFunc("/api/bankaccount", bankaccountController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/bankaccount/{id}", bankaccountController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/bankaccount", bankaccountController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/bankaccount", bankaccountController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/bankaccount/{id}", bankaccountController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ region
	regionController := regionController.Controller{}
	router.HandleFunc("/api/region", regionController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/region/{id}", regionController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/region", regionController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/region", regionController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/region/{id}", regionController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ port
	portController := portController.Controller{}
	router.HandleFunc("/api/port", portController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/port/{id}", portController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/port", portController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/port", portController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/port/{id}", portController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ category
	categoryController := categoryController.Controller{}
	router.HandleFunc("/api/category", categoryController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/category/{id}", categoryController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/category", categoryController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/category", categoryController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/category/{id}", categoryController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ order form 合同的通用属性
	orderformController := orderformController.Controller{}
	router.HandleFunc("/api/orderform_dropDown", orderformController.GetItems_DropDown(mydb)).Methods("GET")

	// ------------ express order 快递单(独立的，不属于合同)
	expressOrderController := expressOrderController.Controller{}
	router.HandleFunc("/api/expressOrder", expressOrderController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/expressOrder/{id}", expressOrderController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/expressOrder", expressOrderController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/expressOrder", expressOrderController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/expressOrder/{id}", expressOrderController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ sell contract
	sellContractController := sellContractController.Controller{}
	router.HandleFunc("/api/sellcontract", sellContractController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/sellcontract/{id}", sellContractController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/sellcontract", sellContractController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/sellcontract", sellContractController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/sellcontract/{id}", sellContractController.DeleteItem(mydb)).Methods("DELETE")
	router.HandleFunc("/api/sellcontract_print/{id}/{templateFolder}/{template}/{printFormat}", sellContractController.Print(mydb)).Methods("GET")

	// customized
	router.HandleFunc("/api/sellcontract_getlast", sellContractController.GetLast(mydb)).Methods("GET")

	// ------------ sell subitem
	sellSubitemController := sellSubitemController.Controller{}
	router.HandleFunc("/api/sellsubitem", sellSubitemController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/sellsubitem/{id}", sellSubitemController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/sellsubitem", sellSubitemController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/sellsubitem", sellSubitemController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/sellsubitem/{id}", sellSubitemController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ buy contract
	buyContractController := buyContractController.Controller{}
	router.HandleFunc("/api/buycontract", buyContractController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/buycontract/{id}", buyContractController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/buycontract", buyContractController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/buycontract", buyContractController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/buycontract/{id}", buyContractController.DeleteItem(mydb)).Methods("DELETE")
	router.HandleFunc("/api/buycontract_print/{id}/{templateFolder}/{template}/{printFormat}", buyContractController.Print(mydb)).Methods("GET")

	// customized
	router.HandleFunc("/api/buycontract_getlast", buyContractController.GetLast(mydb)).Methods("GET")

	// ------------ buy subitem
	buySubitemController := buySubitemController.Controller{}
	router.HandleFunc("/api/buysubitem", buySubitemController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/buysubitem/{id}", buySubitemController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/buysubitem", buySubitemController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/buysubitem", buySubitemController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/buysubitem/{id}", buySubitemController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ mould contract
	mouldContractController := mouldContractController.Controller{}
	router.HandleFunc("/api/mouldcontract", mouldContractController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/mouldcontract/{id}", mouldContractController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/mouldcontract", mouldContractController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/mouldcontract", mouldContractController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/mouldcontract/{id}", mouldContractController.DeleteItem(mydb)).Methods("DELETE")
	router.HandleFunc("/api/mouldcontract_print/{id}/{templateFolder}/{template}/{printFormat}", mouldContractController.Print(mydb)).Methods("GET")

	// customized
	router.HandleFunc("/api/mouldcontract_getlast", mouldContractController.GetLast(mydb)).Methods("GET")

	// ------------ product
	productController := productController.Controller{}
	router.HandleFunc("/api/product", productController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/product/{id}", productController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/product", productController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/product", productController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/product/{id}", productController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/product_dropDown", productController.GetItems_DropDown(mydb)).Methods("GET")
	router.HandleFunc("/api/product_dropDown_sellContract", productController.GetItems_DropDown_sellContract(mydb)).Methods("GET")
	router.HandleFunc("/api/product_dropDown_sellSubitem", productController.GetItems_DropDown_sellSubitem(mydb)).Methods("GET")

	router.HandleFunc("/api/product_component", productController.GetComponents(mydb)).Methods("GET")
	router.HandleFunc("/api/product_component/{parent_id}/{child_id}", productController.Assemble(mydb)).Methods("POST")
	router.HandleFunc("/api/product_component/{parent_id}/{child_id}", productController.Disassemble(mydb)).Methods("DELETE")

	// commodity_product
	router.HandleFunc("/api/commodity_getproduct", productController.GetItems_ByCommodity(mydb)).Methods("GET")

	// ------------ product purchase
	productPurchaseController := productPurchaseController.Controller{}
	router.HandleFunc("/api/productPurchase", productPurchaseController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/productPurchase/{id}", productPurchaseController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/productPurchase", productPurchaseController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/productPurchase", productPurchaseController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/productPurchase/{id}", productPurchaseController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/productPurchase_companySearch", productPurchaseController.GetItems_GroupByCompany(mydb)).Methods("GET")
	router.HandleFunc("/api/productPurchase_historySearch", productPurchaseController.GetItems_History(mydb)).Methods("GET")
	router.HandleFunc("/api/productPurchase_byProductId/{id}/{company_id}", productPurchaseController.GetItem_ByProductId(mydb)).Methods("GET")

	// ------------ commodity
	commodityController := commodityController.Controller{}
	router.HandleFunc("/api/commodity", commodityController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/commodity/{commodity_id}", commodityController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/commodity/{commodity_id}/{product_id}", commodityController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/commodity", commodityController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/commodity", commodityController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/commodity/{id}", commodityController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/commodity_byproduct", commodityController.GetItems_ByProduct(mydb)).Methods("GET")
	// commodity_product
	router.HandleFunc("/api/commodity_byproduct/{commodity_id}/{product_id}", commodityController.Assemble(mydb)).Methods("POST")
	router.HandleFunc("/api/commodity_byproduct/{commodity_id}/{product_id}", commodityController.Disassemble(mydb)).Methods("DELETE")

	// ------------ paymentRequest 内部财务用的账户
	paymentRequestController := paymentRequestController.Controller{}
	router.HandleFunc("/api/paymentRequest", paymentRequestController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/paymentRequest/{id}", paymentRequestController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/paymentRequest", paymentRequestController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/paymentRequest", paymentRequestController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/paymentRequest/{id}", paymentRequestController.DeleteItem(mydb)).Methods("DELETE")

	router.HandleFunc("/api/paymentRequest_approve", paymentRequestController.UpdateItem_approve(mydb)).Methods("PUT")
	router.HandleFunc("/api/paymentRequest_reject", paymentRequestController.UpdateItem_reject(mydb)).Methods("PUT")

	// ------------ financialAccount 内部财务用的账户
	financialaccountController := financialaccountController.Controller{}
	router.HandleFunc("/api/financialAccount", financialaccountController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/financialAccount/{id}", financialaccountController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/financialAccount", financialaccountController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/financialAccount", financialaccountController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/financialAccount/{id}", financialaccountController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ financialLedger 科目树
	financialledgerController := financialledgerController.Controller{}
	router.HandleFunc("/api/financialLedger", financialledgerController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/financialLedger/{id}", financialledgerController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/financialLedger", financialledgerController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/financialLedger", financialledgerController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/financialLedger/{id}", financialledgerController.DeleteItem(mydb)).Methods("DELETE")

	// ------------ financialTransaction 交易明细
	financialTransactionController := financialTransactionController.Controller{}
	router.HandleFunc("/api/financialTransaction", financialTransactionController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/financialTransaction/{id}", financialTransactionController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/financialTransaction", financialTransactionController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/financialTransaction", financialTransactionController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/financialTransaction/{id}", financialTransactionController.DeleteItem(mydb)).Methods("DELETE")
	router.HandleFunc("/api/financialTransaction_print/list/{templateFolder}/{template}/{printFormat}", financialTransactionController.PrintList(mydb)).Methods("GET")

	// ------------ financialVoucher 交易凭证
	financialVoucherController := financialVoucherController.Controller{}
	router.HandleFunc("/api/financialVoucher", financialVoucherController.GetItems(mydb)).Methods("GET")
	router.HandleFunc("/api/financialVoucher/{id}", financialVoucherController.GetItem(mydb)).Methods("GET")
	router.HandleFunc("/api/financialVoucher", financialVoucherController.AddItem(mydb)).Methods("POST")
	router.HandleFunc("/api/financialVoucher", financialVoucherController.UpdateItem(mydb)).Methods("PUT")
	router.HandleFunc("/api/financialVoucher/{id}", financialVoucherController.DeleteItem(mydb)).Methods("DELETE")
	router.HandleFunc("/api/financialVoucher_print/list/{templateFolder}/{template}/{printFormat}", financialVoucherController.PrintList(mydb)).Methods("GET")

}
