package financialTransactionController

import (
	"net/http"
	"reflect"

	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/financialTransaction"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.FinancialTransaction

var authName = "financial"

// =============================================== HTTP REQUESTS
func (c Controller) GetItems(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, mydb)
	}
}

func (c Controller) GetItems_DropDown(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems_DropDown(w, r, mydb)
	}
}

func (c Controller) GetItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItem(w, r, mydb)
	}
}
func (c Controller) AddItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_AddItem(w, r, mydb)
	}
}

func (c Controller) UpdateItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_UpdateItem(w, r, mydb)
	}
}

func (c Controller) DeleteItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_DeleteItem(w, r, mydb)
	}
}

func (c Controller) PrintList(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_PrintList(w, r, mydb)
	}
}

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItems_DropDown(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(mydb, w, r, reflect.TypeOf(item), "GetRows_DropDown", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	if err != nil {
		utils.Log(err, "事务运行失败")
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	if err != nil {
		utils.Log(err, "事务运行失败")
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	if err != nil {
		utils.Log(err, "事务运行失败")
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_PrintList(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	// 从param里取出模板所在目录，打印格式（做在utils里面是为了方便日后修改）
	_, _path, printFormat := utils.FetchPrintPathAndId(r)

	// 读取行数据(和上面搜索的代码一样)
	repo := repository.Repository{}

	// 生成打印数据(取map出来而不是item，是为了方便篡改)
	dataSource, err := repo.GetPrintSourceList(mydb, r, userId)

	if err != nil {
		w.Write([]byte("error on generating source data," + err.Error()))
	}

	// 直接打印到writer(因为打印完毕需要删除cache，所以要在删除之前使用writer)
	err = utils.PrintFromTemplate(w, dataSource, _path, printFormat, userId)

	if err != nil {
		w.Write([]byte("error on printing," + err.Error()))
		return
	}
}

// func (c Controller) Print(mydb models.MyDb) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		pass, userId := auth.CheckAuth(mydb, w, r, authName)
// 		if !pass {
// 			return
// 		}

// 		// 从param里取出id，模板所在目录，打印格式（做在utils里面是为了方便日后修改）
// 		id, _path, printFormat := utils.FetchPrintPathAndId(r)

// 		// 生成打印数据(取map出来而不是item，是为了方便篡改)
// 		repo := repository.Repository{}
// 		dataSource, err := repo.GetPrintSource(mydb, id, userId)

// 		if err != nil {
// 			w.Write([]byte("error on generating source data," + err.Error()))
// 		}

// 		// 直接打印到writer(因为打印完毕需要删除cache，所以要在删除之前使用writer)
// 		err = utils.PrintFromTemplate(w, dataSource, _path, printFormat, userId)

// 		if err != nil {
// 			w.Write([]byte("error on printing," + err.Error()))
// 			return
// 		}
// 	}
// }
