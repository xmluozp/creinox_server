package folderContrller

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/xmluozp/creinox_server/auth"
	imageController "github.com/xmluozp/creinox_server/controllers/imagedata"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/folder"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Folder

var authName = ""

// =============================================== HTTP REQUESTS

func (c Controller) AddItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_AddItem(w, r, mydb)
	}
}

func (c Controller) Print(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_Print(w, r, mydb)
	}
}

// =============================================== basic CRUD
func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
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

// ============================== internal

// 删除folder把对应的image删掉。id 是 folder_id
func (c Controller) Delete(mydb models.MyDb, folder_id int, userId int) error {

	imageCtrl := imageController.Controller{}

	repo := repository.Repository{}
	_, err := repo.DeleteRow(mydb, folder_id, userId)

	// folder, err := repo.GetRow(mydb, folder_id, userId)

	// if err != nil {
	// 	fmt.Println("删除folder时取folder失败")
	// 	utils.Log(err)
	// 	return err
	// }

	images, err := imageCtrl.ItemsByFolder(mydb, folder_id, userId)

	if err != nil {
		fmt.Println("删除folder时取下属images失败")
		utils.Log(err)
		return err
	}

	// fmt.Println("images", images)

	for key := range images {
		err = imageCtrl.Delete(mydb, images[key].ID.Int, userId)
		if err != nil {
			fmt.Println("删除folder时删除下属images失败")
			utils.Log(err)
			return err
		}
	}

	return err
}
