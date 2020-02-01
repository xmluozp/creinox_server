package imageControllr

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/gobuffalo/nulls"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/imagedata"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Image

var authName = ""
var UPLOAD_FOLDER = "uploads/"

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

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, _ := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo)
		utils.SendJson(w, status, returnValue, err)
	}
}
func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
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

// 这里是前端的文件夹管理用，特殊处理
// func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {

// 	return func(w http.ResponseWriter, r *http.Request) {

// 		pass, userId := auth.CheckAuth(db, w, r, authName)
// 		if !pass {
// 			return
// 		}

// 		var item modelName
// 		repo := repository.Repository{}

// 		status, returnValue, deletedItem, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)

// 		// 取出旧图片删掉
// 		imageOld := deletedItem.(models.Image)
// 		os.Remove(UPLOAD_FOLDER + imageOld.Path.String)
// 		os.Remove(UPLOAD_FOLDER + imageOld.ThumbnailPath.String)

// 		utils.SendJson(w, status, returnValue, err)
// 	}
// }

// 这里是前端的文件夹管理用，特殊处理：手动读取和删除
func (c Controller) DeleteItems(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var returnValue models.JsonRowsReturn
		status := http.StatusOK
		var idList []int
		err := json.NewDecoder(r.Body).Decode(&idList)

		for _, value := range idList {

			err = c.Delete(db, value, userId)
			if err != nil {
				returnValue.Info += fmt.Sprintf("删除图片出错，ID：%d, %s", value, err.Error())
				status = http.StatusBadRequest
			}
		}

		if err == nil { // no err
			returnValue.Info = fmt.Sprintf("删除%d张图片", len(idList))
		}

		fmt.Println("删完了信息有问题：", returnValue.Info, returnValue, err)
		utils.SendJson(w, status, returnValue, err)
	}
}

// 批量上传
func (c Controller) AddItems(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		// ---------------
		// 取出提交的files和folder_id
		params := mux.Vars(r)
		folder_id, err := strconv.Atoi(params["folder_id"])

		if err != nil {
			fmt.Println("test: folder_id get error", err)
			return
		}

		status, returnValue, _, files, err := utils.DecodeFormData(r, reflect.TypeOf(item))

		if err != nil {
			utils.SendJson(w, status, returnValue, err)
			return
		}

		// 循环存入images
		for key := range files {
			c.Upload(db, -1, key, files[key], folder_id, userId) // 不需要删旧数据也不需要返回，所以直接存就好
		}

		returnValue.Info = fmt.Sprintf("上传了%d张图片", len(files))
		// ---------------
		utils.SendJson(w, status, returnValue, err)
	}
}
func (c Controller) Show(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		file, _ := os.Open("." + r.URL.Path)
		// errorHandle(err, w);

		defer file.Close()
		buff, _ := ioutil.ReadAll(file)
		// errorHandle(err, w);
		w.Write(buff)
	}
}

// TODO: 批量删除：比如公司下的文件夹管理会用到

// ============================== internal

func (c Controller) ItemsByFolder(
	db *sql.DB,
	folderId int) ([]modelName, error) {

	repo := repository.Repository{}
	images, err := repo.GetRowsByFolder(db, folderId)

	return images, err
}
func (c Controller) Item(db *sql.DB, id int) (image modelName, err error) {

	repo := repository.Repository{}
	image, err = repo.GetRow(db, id)

	return image, err
}

// input: old table, old column, old id, the file. output: new imagedata id
func (c Controller) Upload(
	db *sql.DB,
	oldImage_id int,
	fileName string,
	fileBytes []byte,
	folder_id int,
	userId int) (int, error) {

	// 先upload file

	// -------------- upload original file
	tempFile, err := ioutil.TempFile("uploads", "pic_*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	defer os.Remove(tempFile.Name())

	// -------------- thumbnail
	imagePath, _ := os.Open(tempFile.Name())
	defer imagePath.Close()
	srcImage, _, err := image.Decode(imagePath)

	if err != nil {
		fmt.Println(err)
	}

	dx := srcImage.Bounds().Dx()
	dy := srcImage.Bounds().Dy()

	// Dimension of new thumbnail 80 X 80
	dstImage := image.NewRGBA(image.Rect(0, 0, 200, 200*dy/dx))

	// Thumbnail function of Graphics
	err = graphics.Thumbnail(dstImage, srcImage)

	if err != nil {
		return 0, err
	}

	fileInfo, err := os.Stat(tempFile.Name())

	if err != nil {
		return 0, err
	}

	fmt.Println("os:", fileInfo.Name())

	newImage, err := os.Create("uploads/thumbnail_" + fileInfo.Name())

	if err != nil {
		return 0, err
	}

	// defer os.Remove("uploads/thumbnail_" + fileInfo.Name())

	defer newImage.Close()

	png.Encode(newImage, dstImage)

	// 再执行SQL, 生成image record（upload file出错概率高，先运行，如果panic，就从defer地方delete
	var newImagedata models.Image

	newImagedata.Name = nulls.NewString(fileName)
	newImagedata.Height = nulls.NewInt(dy)
	newImagedata.Width = nulls.NewInt(dx)
	newImagedata.Path = nulls.NewString(fileInfo.Name())
	newImagedata.ThumbnailPath = nulls.NewString("thumbnail_" + fileInfo.Name())
	newImagedata.Ext = nulls.NewString(filepath.Ext(fileInfo.Name()))

	if folder_id > 0 {
		newImagedata.Gallary_folder_id = nulls.NewInt(folder_id)
	}

	repo := repository.Repository{}

	newImagedataResult, err := repo.AddRow(db, newImagedata, userId)

	if err != nil {
		return 0, err
	}

	// 删除原图：
	if oldImage_id > 0 {
		c.Delete(db, oldImage_id, userId)
	}

	// 忽略错误。因为有可能数据库没图片
	return newImagedataResult.ID, nil
}

func (c Controller) Delete(
	db *sql.DB,
	id int,
	userId int) error {

	repo := repository.Repository{}

	// 删除原图：
	deletedItem, err := repo.DeleteRow(db, id, userId)

	// 如果没有原图就不管
	if err != nil {
		// fmt.Println("删除数据库图片记录出错", err.Error())
		return nil
	}

	imageOld := deletedItem.(models.Image)

	os.Remove(UPLOAD_FOLDER + imageOld.Path.String)
	os.Remove(UPLOAD_FOLDER + imageOld.ThumbnailPath.String)
	return nil
}
