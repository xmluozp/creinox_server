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
	"github.com/Unknwon/goconfig"
	"github.com/gobuffalo/nulls"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repositoryFolder "github.com/xmluozp/creinox_server/repository/folder"
	repository "github.com/xmluozp/creinox_server/repository/imagedata"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Image

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

// 通用图片管理功能用。否则只有批量删除
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

		// 取图片列表。假如没有folder，前端会传folder_structure，这里顺便也取出来
		var item models.Folder
		status, returnValue, folderItem, files, err := utils.DecodeFormData(r, reflect.TypeOf(item))

		// ---------------
		// 取出提交的files和folder_id
		params := mux.Vars(r)
		folder_id, err := strconv.Atoi(params["folder_id"])

		// 原数据库的公司是没有folder的

		if err != nil {
			utils.SendJson(w, status, returnValue, err)
			return
		}

		// folder id -1说明需要生成folder
		if folder_id == -1 {

			folderRepo := repositoryFolder.Repository{}

			newFolder, err := folderRepo.AddRow_withRef(db, folderItem.(models.Folder), userId)

			folder_id = newFolder.ID.Int

			if err != nil {

				fmt.Println("folder生成出错", err)
				utils.SendJson(w, status, returnValue, err)
				return
			}
		}

		// 循环存入images
		for key := range files {
			c.Upload(db, -1, key, files[key], folder_id, userId) // 不需要删旧数据也不需要返回，所以直接存就好
		}

		returnValue.Info = fmt.Sprintf("上传了%d张图片", len(files))

		// 把folder的id回传给前端（因为假如folder不存在，有可能新建folder）
		returnValue.Row = models.Folder{ID: nulls.NewInt(folder_id)}

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

// ============================== internal

// Get item by folder; Folder will be connected with company/product/other tables
func (c Controller) ItemsByFolder(
	db *sql.DB,
	folderId int,
	userId int) ([]modelName, error) {

	repo := repository.Repository{}
	images, err := repo.GetRowsByFolder(db, folderId, userId)

	return images, err
}

// Get item. Will be called by other controller
func (c Controller) Item(db *sql.DB, id int) (image modelName, err error) {

	repo := repository.Repository{}
	image, err = repo.GetRow(db, id, 0)

	return image, err
}

// Upload withouth a front-end-expected json return. input: old table, old column, old id, the file. output: new imagedata id
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
	return newImagedataResult.ID.Int, nil
}

// Delete withouth a front-end-expected json return
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

	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}

	uploads, err := cfg.GetValue("site", "uploads")

	os.Remove(uploads + "/" + imageOld.Path.String)
	os.Remove(uploads + "/" + imageOld.ThumbnailPath.String)
	return nil
}
