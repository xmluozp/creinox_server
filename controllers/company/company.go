package companyController

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"net/http"
	"strconv"
	"strings"

	"reflect"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/auth"
	folderController "github.com/xmluozp/creinox_server/controllers/folder"
	imageController "github.com/xmluozp/creinox_server/controllers/imagedata"

	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/company"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Company

var authNames = []string{
	"",
	"companyinternal",
	"companyfactory",
	"companyoverseas",
	"companydomestic",
	"companyshipping"}

// =============================================== basic CRUD

func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		searchTerms := utils.GetSearchTerms(r)
		companyType, _ := strconv.Atoi(searchTerms["companyType"])

		pass, userId := auth.CheckAuth(db, w, r, authNames[companyType])
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 数据权限没做
		pass, userId := auth.CheckAuth(db, w, r, "")
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

		itemFromRequest := returnItem.(models.Company)

		// 更新image数据库, 上传图片
		err = updateImage(db, itemFromRequest, files, userId)

		if err != nil {
			var returnValue models.JsonRowsReturn
			returnValue.Info = "文件上传错误" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

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

		fmt.Println("看看营业执照：", itemFromRequest.ImageLicense_id)

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
		company := itemReturn.(modelName)

		if err == nil {

			// 删除营业执照和名片
			imageCtrl := imageController.Controller{}
			imageCtrl.Delete(db, company.ImageLicense_id.Int, userId)
			imageCtrl.Delete(db, company.ImageBizCard_id.Int, userId)

			// 删除folder (folder下面images的删除在folder处理)
			folderCtrl := folderController.Controller{}
			err = folderCtrl.Delete(db, company.Gallary_folder_id.Int, userId)

			if err != nil {
				returnValue.Info = "删除公司对应图库失败" + err.Error()
				utils.SendJson(w, http.StatusFailedDependency, returnValue, err)
			}
		}

		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) Print(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, "")
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

// 以下这段代码找不到方法generalization，只好复制粘贴到各自的controllers里。但代码都是一样的：把业务表里附带的image update到image表、上传到文件夹、更新业务表对应fk
func updateImage(db *sql.DB, item modelName, files map[string][]byte, userId int) error {

	repo := repository.Repository{}

	// update
	updatedItem, _ := repo.GetRow(db, item.ID.Int, 0)
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
		fileName := fmt.Sprintf("company.image_id.%d", item.ID.Int)
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
