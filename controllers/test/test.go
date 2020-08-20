package testController

import (
	"database/sql"
	"fmt"
	"net/http"

	// repository "github.com/xmluozp/creinox_server/repository/bankAccount"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}

// type modelName = models.BankAccount

var authName = ""

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func (c Controller) C_GetItems_DropDown(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}
func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) Test(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		v, _ := params["v"]

		a1 := utils.ParseFlightSlice(v)

		fmt.Println("Test output:", a1)
		// returnValue.Message =
		// utils.SendJson(w, http.StatusOK, "hi", nil)

	}
}

func (c Controller) TestApp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		v, _ := params["v"]

		a1 := utils.ParseFlightSlice(v)

		fmt.Println("Test output:", a1)
		// returnValue.Message =
		// utils.SendJson(w, http.StatusOK, "hi", nil)

	}
}
