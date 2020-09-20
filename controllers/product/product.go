package productController

import (
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

func (c Controller) GetItems_DropDown_sellContract(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems_DropDown_sellContract(w, r, mydb)
	}
}

func (c Controller) GetItems_DropDown_sellSubitem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems_DropDown_sellSubitem(w, r, mydb)
	}
}

func (c Controller) GetItems_ByCommodity(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems_ByCommodity(w, r, mydb)
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

func (c Controller) C_GetItems_DropDown_sellContract(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(mydb, w, r, reflect.TypeOf(item), "GetRows_DropDown_sellContract", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItems_DropDown_sellSubitem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(mydb, w, r, reflect.TypeOf(item), "GetRows_DropDown_sellSubitem", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItems_ByCommodity(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(mydb, w, r, reflect.TypeOf(item), "GetRows_ByCommodity", repo, userId)
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

	pass, userId := auth.CheckAuth(mydb, w, r, "")
	if !pass {
		return
	}

	var item modelName
	var returnValue models.JsonRowsReturn

	repo := repository.Repository{}
	// f, _, _ := utils.GetFunc_AddWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	defer tx.Rollback()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, returnItem, files, err := utils.GetFunc_AddWithHTTPReturn_FormData(mydb, w, r, reflect.TypeOf(item), repo, userId)

	// 验证不通过之类的问题就不需要传图
	if err != nil {
		utils.SendJson(w, status, returnValue, err)
		return
	}

	itemFromRequest := returnItem.(modelName)

	// 更新category里的最大编码
	ca := categoryController.Controller{}
	// _, currentCode := utils.ParseFlight(itemFromRequest.Code.String)
	// currentCodeSlice := utils.ParseFlightSlice(itemFromRequest.Code.String)
	// var currentCode string
	// if len(currentCodeSlice) > 2 {
	// 	currentCode = currentCodeSlice[1]
	// } else {
	// 	currentCode = ""
	// }

	err = ca.Update_currentCode(mydb, itemFromRequest.Category_id.Int, itemFromRequest.Code.String, userId)

	if err != nil {
		returnValue.Info = "编码更新错误" + err.Error()
		utils.SendError(w, http.StatusInternalServerError, returnValue)
		return
	}

	// 更新image数据库, 上传图片
	err = updateImage(mydb, itemFromRequest, files, userId)

	if err != nil {
		returnValue.Info = "文件上传错误" + err.Error()
		utils.SendError(w, http.StatusInternalServerError, returnValue)
		return
	}

	// 假如设置为商品，就更新商品表
	if itemFromRequest.IsCreateCommodity.Bool {
		commodityCtrl := commodityController.Controller{}
		err = commodityCtrl.Add_ByProduct(mydb, itemFromRequest.ID.Int, userId)

		if err != nil {
			returnValue.Info = "设置为商品时出错" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}
	}

	err = tx.Commit()

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, "")
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	defer tx.Rollback()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	// upload form
	status, returnValue, returnItem, files, err := utils.GetFunc_UpdateWithHTTPReturn_FormData(mydb, w, r, reflect.TypeOf(item), repo, userId)

	// 验证不通过之类的问题就不需要传图
	if err != nil {
		utils.SendJson(w, status, returnValue, err)
		return
	}

	// convert "reflected" item into company type
	itemFromRequest := returnItem.(modelName)

	// 更新产品图片图片. 如果没有就是删除
	err = updateImage(mydb, itemFromRequest, files, userId)

	if err != nil {
		var returnValue models.JsonRowsReturn
		returnValue.Info = "文件上传错误" + err.Error()
		utils.SendError(w, http.StatusInternalServerError, returnValue)
		return
	}

	// 取最新row返回
	updatedItem, err := repo.GetRow(mydb, itemFromRequest.ID.Int, userId)
	returnValue.Row = updatedItem

	err = tx.Commit()

	// send success message to front-end
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, "")
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	defer tx.Rollback()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, itemReturn, err := utils.GetFunc_DeleteWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	product := itemReturn.(modelName)

	if err == nil {

		// 删除产品图片
		imageCtrl := imageController.Controller{}
		err = imageCtrl.Delete(mydb, product.Image_id.Int, userId)

		if err != nil {
			utils.Log(err, "删除产品对应图库失败")
			returnValue.Info = "删除产品对应图库失败, " + err.Error()
			utils.SendJson(w, http.StatusFailedDependency, returnValue, err)
			return
		}
	}

	err = tx.Commit()
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

// ================= components
func (c Controller) GetComponents(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(mydb, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(mydb, w, r, reflect.TypeOf(item), "GetRows_Component", repo, userId)

		fmt.Println("*************")
		fmt.Println(w, status, returnValue, err)
		fmt.Println("*************")

		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) Assemble(mydb models.MyDb) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(mydb, w, r, authName)
		if !pass {
			return
		}

		repo := repository.Repository{}

		// -------------- customized
		var returnValue models.JsonRowsReturn
		params := mux.Vars(r)

		parent_id, _ := strconv.Atoi(params["parent_id"])
		child_id, _ := strconv.Atoi(params["child_id"])

		err := repo.Assemble(mydb, parent_id, child_id, userId)

		if err != nil {
			returnValue.Info = err.Error()
			utils.SendJson(w, http.StatusBadRequest, returnValue, err)
			return
		}

		utils.SendJson(w, http.StatusOK, returnValue, err)
	}
}

func (c Controller) Disassemble(mydb models.MyDb) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(mydb, w, r, authName)
		if !pass {
			return
		}

		repo := repository.Repository{}

		// -------------- customized
		var returnValue models.JsonRowsReturn
		params := mux.Vars(r)

		parent_id, _ := strconv.Atoi(params["parent_id"])
		child_id, _ := strconv.Atoi(params["child_id"])

		err := repo.Disassemble(mydb, parent_id, child_id, userId)

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
func updateImage(mydb models.MyDb, item modelName, files map[string][]byte, userId int) error {

	repo := repository.Repository{}

	// update
	updatedItem, _ := repo.GetRow(mydb, item.ID.Int, userId)
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
		fileName := fmt.Sprintf("product.%s.%d", columnName, item.ID.Int)
		newImageId, err := imageCtrl.Upload(mydb, oldImage_id, fileName, fileBytes, -1, userId)

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
		_, err = repo.UpdateRow(mydb, updatedItem, userId)

		// send error of upload file to front-end
		if err != nil {
			return err
		}
	}

	return nil
}
