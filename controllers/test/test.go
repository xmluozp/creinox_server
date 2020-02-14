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

func (c Controller) Test(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		v, _ := params["v"]

		a1, a2 := utils.ParseFlight(v)

		fmt.Println("Test output:", a1, a2)
		// returnValue.Message =
		// utils.SendJson(w, http.StatusOK, "hi", nil)

	}
}
