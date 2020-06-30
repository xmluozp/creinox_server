package productController

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/nulls"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	categoryController "github.com/xmluozp/creinox_server/controllers/category"
	commodityController "github.com/xmluozp/creinox_server/controllers/commodity"
	imageController "github.com/xmluozp/creinox_server/controllers/imagedata"

	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/product"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Product

var authName = "product"

// =============================================== basic CRUD

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

func (c Controller) GetItems_DropDown_sellContract(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_DropDown_sellContract", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItems_DropDown_sellSubitem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_DropDown_sellSubitem", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItems_ByCommodity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_ByCommodity", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}
func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, "")
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		// f, _, _ := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)

		status, returnValue, returnItem, files, err := utils.GetFunc_AddWithHTTPReturn_FormData(db, w, r, reflect.TypeOf(item), repo, userId)

		// 验证不通过之类的问题就不需要传图
		if err != nil {
			utils.SendJson(w, status, returnValue, err)
			return
		}

		itemFromRequest := returnItem.(modelName)

		// 更新category里的最大编码
		ca := categoryController.Controller{}
		_, currentCode := utils.ParseFlight(itemFromRequest.Code.String)
		ca.Update_currentCode(db, itemFromRequest.Category_id.Int, currentCode, userId)

		if err != nil {
			var returnValue models.JsonRowsReturn
			returnValue.Info = "编码更新错误" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		// 更新image数据库, 上传图片
		err = updateImage(db, itemFromRequest, files, userId)

		if err != nil {
			var returnValue models.JsonRowsReturn
			returnValue.Info = "文件上传错误" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		if itemFromRequest.IsCreateCommodity.Bool {
			commodityCtrl := commodityController.Controller{}
			err = commodityCtrl.Add_ByProduct(db, itemFromRequest.ID.Int, userId)

			if err != nil {
				var returnValue models.JsonRowsReturn
				returnValue.Info = "设置为商品时出错" + err.Error()
				utils.SendError(w, http.StatusInternalServerError, returnValue)
				return
			}
		}
		// 假如设置为商品，就更新商品表

		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) UpdateItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, "")
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		// upload form
		status, returnValue, returnItem, files, err := utils.GetFunc_UpdateWithHTTPReturn_FormData(db, w, r, reflect.TypeOf(item), repo, userId)

		// 验证不通过之类的问题就不需要传图
		if err != nil {
			utils.SendJson(w, status, returnValue, err)
			return
		}

		// convert "reflected" item into company type
		itemFromRequest := returnItem.(modelName)

		// 更新公司的两张图片. 如果没有就是删除
		err = updateImage(db, itemFromRequest, files, userId)

		if err != nil {
			var returnValue models.JsonRowsReturn
			returnValue.Info = "文件上传错误" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		// 取最新row返回
		updatedItem, err := repo.GetRow(db, itemFromRequest.ID.Int, userId)
		returnValue.Row = updatedItem

		// send success message to front-end
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, "")
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, itemReturn, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		product := itemReturn.(modelName)

		if err == nil {

			// 删除产品图片
			imageCtrl := imageController.Controller{}
			err = imageCtrl.Delete(db, product.Image_id.Int, userId)

			if err != nil {
				returnValue.Info = "删除产品对应图库失败" + err.Error()
				utils.SendJson(w, http.StatusFailedDependency, returnValue, err)
			}
		}

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

// ================= components
func (c Controller) GetComponents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRows_Component", repo, userId)

		fmt.Println("*************")
		fmt.Println(w, status, returnValue, err)
		fmt.Println("*************")

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

		parent_id, _ := strconv.Atoi(params["parent_id"])
		child_id, _ := strconv.Atoi(params["child_id"])

		err := repo.Assemble(db, parent_id, child_id, userId)

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

		parent_id, _ := strconv.Atoi(params["parent_id"])
		child_id, _ := strconv.Atoi(params["child_id"])

		err := repo.Disassemble(db, parent_id, child_id, userId)

		if err != nil {
			returnValue.Info = err.Error()
			utils.SendJson(w, http.StatusBadRequest, returnValue, err)
			return
		}

		utils.SendJson(w, http.StatusOK, returnValue, err)
	}
}

// 以下这段代码找不到方法generalization，只好复制粘贴到各自的controllers里。但代码都是一样的：把业务表里附带的image update到image表、上传到文件夹、更新业务表对应fk
// =====image
func updateImage(db *sql.DB, item modelName, files map[string][]byte, userId int) error {

	repo := repository.Repository{}

	// update
	updatedItem, _ := repo.GetRow(db, item.ID.Int, userId)
	newImageIds := make(map[string]int)

	// upload image(here will be twice: license, biscard), return new image id -------------------------------------
	for key, _ := range files {

		fileBytes := files[key]

		imageCtrl := imageController.Controller{}

		// magical convention: 截取"."之前的字符当做column name (因为update的时候需要删掉原图，所以需要这个colunm)
		// 截它是因为它的格式是image_id.row
		columnName := strings.Split(key, ".")[0]

		oldImage_id := -1
		oldImage_id_value := utils.GetFieldValue(columnName, "json", updatedItem)

		if oldImage_id_value != nil {
			oldImage_id = oldImage_id_value.(nulls.Int).Int
		}

		var err error
		fileName := fmt.Sprintf("product.image_id.%d", item.ID.Int)
		newImageId, err := imageCtrl.Upload(db, oldImage_id, fileName, fileBytes, -1, userId)

		if newImageId != 0 {
			newImageIds[columnName] = newImageId
		}

		if err != nil {
			return err
		}
	}

	// insert newid back
	if len(newImageIds) > 0 {
		jsonString, err := json.Marshal(newImageIds)
		json.Unmarshal(jsonString, &updatedItem)
		_, err = repo.UpdateRow(db, updatedItem, userId)

		// send error of upload file to front-end
		if err != nil {
			return err
		}
	}

	return nil
}
