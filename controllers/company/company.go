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

		pass, _ := auth.CheckAuth(db, w, r, authNames[companyType])
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var item modelName
		repo := repository.Repository{}
		status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo)
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

		f, returnItem, files, err := utils.GetFunc_AddWithHTTPReturn_FormData(db, w, r, reflect.TypeOf(item), repo, userId)
		if err != nil {
			f()
			fmt.Println("update error", err.Error())
			return
		}

		itemFromRequest := returnItem.(models.Company)

		err = updateImage(db, itemFromRequest, files, userId)

		if err != nil {
			var returnValue models.JsonRowsReturn
			returnValue.Info = "文件上传错误" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		f()
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
		f, returnItem, files, err := utils.GetFunc_UpdateWithHTTPReturn_FormData(db, w, r, reflect.TypeOf(item), repo, userId)
		if err != nil {
			f()
			fmt.Println("update error", err.Error())
			return
		}

		// convert "reflected" item into company type
		itemFromRequest := returnItem.(models.Company)

		err = updateImage(db, itemFromRequest, files, userId)

		if err != nil {
			var returnValue models.JsonRowsReturn
			returnValue.Info = "文件上传错误" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		// send success message to front-end
		f()
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

		f, company, _ := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		f()

		itemReturn := company.(models.Company)

		fmt.Println("返回外层", itemReturn)

	}
}

// ===================================================
func updateImage(db *sql.DB, company models.Company, files map[string][]byte, userId int) error {

	repo := repository.Repository{}

	// update
	updatedItem, _ := repo.GetRow(db, company.ID)
	newImageIds := make(map[string]int)

	// upload image(here will be twice: license, biscard), return new image id -------------------------------------
	for key, _ := range files {

		fileBytes := files[key]

		imageCtrl := imageController.Controller{}

		// magical convention: 截取"."之前的字符当做column name (因为update的时候需要删掉原图，所以需要这个colunm)
		// 截它是因为它的格式是image_id.row
		columnName := strings.Split(key, ".")[0]

		var oldImage_id int

		oldImage_id = utils.GetFieldValue(columnName, "json", updatedItem).(nulls.Int).Int

		var err error
		newImageId, err := imageCtrl.Upload(db, oldImage_id, fileBytes, userId)

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
