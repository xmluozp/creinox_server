package imageControllr

import (
	"database/sql"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/imagedata"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Company

var authName = ""

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
		f, _, _ := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		f()
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
		f, _, _ := utils.GetFunc_UpdateWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		f()
	}
}

// TODO: 删除对应image
func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		f, _, _ := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		f()
	}
}

// TODO: 批量删除

// ============================== internal interfaces

// input: old table, old column, old id, the file. output: new imagedata id
func (c Controller) Upload(
	db *sql.DB,
	oldImage_id int,
	fileBytes []byte,
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
	graphics.Thumbnail(dstImage, srcImage)

	fileInfo, err := os.Stat(tempFile.Name())

	if err != nil {
		return 0, err
	}

	fmt.Println("os:", fileInfo.Name())

	newImage, err := os.Create("uploads/thumbnail_" + fileInfo.Name())

	defer os.Remove("uploads/thumbnail_" + fileInfo.Name())

	defer newImage.Close()

	png.Encode(newImage, dstImage)

	// 再执行SQL, 生成image record（upload file出错概率高，先运行，如果panic，就从defer地方delete
	var newImagedata models.Image

	newImagedata.Name = nulls.NewString(tempFile.Name())
	newImagedata.Height = nulls.NewInt(dy)
	newImagedata.Width = nulls.NewInt(dx)
	newImagedata.Path = nulls.NewString(fileInfo.Name())
	newImagedata.ThumbnailPath = nulls.NewString("thumbnail_" + fileInfo.Name())
	newImagedata.Ext = nulls.NewString(filepath.Ext(fileInfo.Name()))

	repo := repository.Repository{}
	newImagedataResult, err := repo.AddRow(db, newImagedata, userId)

	// 删除原图：

	// 然后删除原image
	fmt.Println("image return id", newImagedataResult.ID)
	return newImagedataResult.ID, nil
}
