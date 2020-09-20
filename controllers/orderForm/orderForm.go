package orderformController

import (
	"net/http"
	"reflect"

	// xlst "github.com/ivahaev/go-xlsx-templater"

	"github.com/xmluozp/creinox_server/auth"

	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/orderForm"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.OrderForm

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
