package orderformController

import (
	"database/sql"
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

func (c Controller) C_GetItems_DropDown(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_DropDown", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

// =============================================== HTTP REQUESTS
func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, db)
	}
}

func (c Controller) GetItems_DropDown(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems_DropDown(w, r, db)
	}
}
