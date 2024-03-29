package mouldContractController

import (
	"net/http"
	"reflect"

	// xlst "github.com/ivahaev/go-xlsx-templater"

	"github.com/xmluozp/creinox_server/auth"
	folderController "github.com/xmluozp/creinox_server/controllers/folder"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/mouldContract"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.MouldContract

var authName = "mouldcontract"

// =============================================== HTTP REQUESTS
func (c Controller) GetItems(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, mydb)
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

func (c Controller) Print(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_Print(w, r, mydb)
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
		tx.Rollback()
	} else {
		tx.Commit()
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
		tx.Rollback()
	} else {
		tx.Commit()
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

	status, returnValue, itemReturn, err := utils.GetFunc_DeleteWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	contract := itemReturn.(modelName)

	if err == nil {
		// 删除folder (folder下面images的删除在folder处理)
		folderCtrl := folderController.Controller{}
		err = folderCtrl.Delete(mydb, contract.Gallary_folder_id.Int, userId)

		if err != nil {
			returnValue.Info = "删除产品开发合同对应图库失败" + err.Error()
			utils.SendJson(w, http.StatusFailedDependency, returnValue, err)
		}
	}

	if err != nil {
		utils.Log(err, "事务运行失败")
		tx.Rollback()
	} else {
		tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	// 从param里取出id，模板所在目录，打印格式（做在utils里面是为了方便日后修改）
	id, _path, printFormat := utils.FetchPrintPathAndId(r)

	// 生成打印数据(取map出来而不是item，是为了方便篡改)
	repo := repository.Repository{}
	dataSource, err := repo.GetPrintSource(mydb, id, userId)

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

// =============================================== customized

func (c Controller) GetLast(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(mydb, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_RowWithHTTPReturn_MethodName(mydb, w, r, reflect.TypeOf(item), "GetRow_GetLast", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}
