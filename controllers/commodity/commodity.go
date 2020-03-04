package commodityController

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/commodity"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Commodity

var authName = "commodity"

// =============================================== basic CRUD

func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, _ := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItems_DropDown(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, _ := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_DropDown", repo)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, _ := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		params := mux.Vars(r)

		commodity_id, err := strconv.Atoi(params["commodity_id"])
		product_id, err := strconv.Atoi(params["product_id"])

		// if commodity_id == 0, get it from product_id
		var item modelName
		var returnValue models.JsonRowsReturn
		repo := repository.Repository{}

		if commodity_id != 0 {
			item, err = repo.GetRow(db, commodity_id)
		} else {
			item, err = repo.GetRow_ByProduct(db, product_id)
		}

		if err != nil {
			// 如果没记录不需要报错，直接返回空
			utils.SendJson(w, http.StatusOK, returnValue, err)
			return
		}

		returnValue.Row = item
		utils.SendJson(w, http.StatusOK, returnValue, err)
	}
}
func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}
		//=================== customized
		var commodity_product models.Commodity_product
		var returnValue models.JsonRowsReturn
		repo := repository.Repository{}

		// 前台传的是中间表。创建时候默认就是isMeta，assemble才是普通的连接
		err := json.NewDecoder(r.Body).Decode(&commodity_product)

		if err != nil {
			returnValue.Info = err.Error()
			utils.SendJson(w, http.StatusBadRequest, returnValue, err)
			return
		}

		_, err = repo.AddRow_WithProduct(db, commodity_product, userId)

		if err != nil {
			returnValue.Info = err.Error()
			utils.SendJson(w, http.StatusBadRequest, returnValue, err)
			return
		}

		utils.SendJson(w, http.StatusOK, returnValue, err)
	}
}

func (c Controller) UpdateItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}

}

// ==================customized
// 用在"新建产品的时候，顺便建立商品"用. provides product_is, create commodity and commodity_product
func (c Controller) Add_ByProduct(db *sql.DB, product_id int, userId int) error {

	var commodity_product models.Commodity_product

	repo := repository.Repository{}

	commodity_product.Product_id = nulls.NewInt(product_id)

	_, err := repo.AddRow_WithProduct(db, commodity_product, userId)

	if err != nil {
		return err
	}

	return nil
}

func (c Controller) GetItems_ByProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, _ := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_ByProduct", repo)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) Assemble(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		repo := repository.Repository{}

		// -------------- customized
		var returnValue models.JsonRowsReturn
		params := mux.Vars(r)

		commodity_id, _ := strconv.Atoi(params["commodity_id"])
		product_id, _ := strconv.Atoi(params["product_id"])

		err := repo.Assemble(db, commodity_id, product_id, userId)

		if err != nil {
			returnValue.Info = err.Error()
			utils.SendJson(w, http.StatusBadRequest, returnValue, err)
			return
		}

		utils.SendJson(w, http.StatusOK, returnValue, err)
	}
}

func (c Controller) Disassemble(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		repo := repository.Repository{}

		// -------------- customized
		var returnValue models.JsonRowsReturn
		params := mux.Vars(r)

		commodity_id, _ := strconv.Atoi(params["commodity_id"])
		product_id, _ := strconv.Atoi(params["product_id"])

		err := repo.Disassemble(db, commodity_id, product_id, userId)

		if err != nil {
			returnValue.Info = err.Error()
			utils.SendJson(w, http.StatusBadRequest, returnValue, err)
			return
		}

		utils.SendJson(w, http.StatusOK, returnValue, err)
	}
}
