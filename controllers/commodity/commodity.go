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

func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItems_DropDown(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_DropDown", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
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
			item, err = repo.GetRow(db, commodity_id, userId)
		} else {
			item, err = repo.GetRow_ByProduct(db, product_id, userId)
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

func (c Controller) Print(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_ByProduct", repo, userId)
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
