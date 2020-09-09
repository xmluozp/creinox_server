package paymentRequestController

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/paymentRequest"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.PaymentRequest

var authName = "paymentRequest"
var authNameConfirm = "confirm-payment"

// =============================================== HTTP REQUESTS
func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, db)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItem(w, r, db)
	}
}
func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_AddItem(w, r, db)
	}
}

func (c Controller) UpdateItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_UpdateItem(w, r, db)
	}
}

func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_DeleteItem(w, r, db)
	}
}

func (c Controller) Print(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_Print(w, r, db)
	}
}

func (c Controller) UpdateItem_approve(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_UpdateItem_approve(w, r, db)
	}
}

func (c Controller) UpdateItem_reject(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_UpdateItem_reject(w, r, db)
	}
}

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	// 从param里取出id，模板所在目录，打印格式（做在utils里面是为了方便日后修改）
	id, _path, printFormat := utils.FetchPrintPathAndId(r)

	// 生成打印数据(取map出来而不是item，是为了方便篡改)
	repo := repository.Repository{}
	dataSource, err := repo.GetPrintSource(db, id, userId)

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

// =============================== customized

func (c Controller) C_UpdateItem_approve(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authNameConfirm)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn_MethodName(db, w, r, reflect.TypeOf(item), "UpdateRow_approve", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem_reject(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authNameConfirm)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn_MethodName(db, w, r, reflect.TypeOf(item), "UpdateRow_reject", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}