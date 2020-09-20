package textTemplateController

import (
	"net/http"
	"reflect"

	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/textTemplate"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.TextTemplate

var authName = ""

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
	status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

// =============================================== customized
func (c Controller) GetItems_Template(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(mydb, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(mydb, w, r, reflect.TypeOf(item), "GetRows_Template", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}
