package folderContrller

import (
	"database/sql"
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

// =============================================== basic CRUD
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

// ============================== internal

// 删除folder把对应的image删掉。id 是 folder_id
func (c Controller) Delete(db *sql.DB, folder_id int, userId int) error {

	imageCtrl := imageController.Controller{}

	repo := repository.Repository{}
	_, err := repo.DeleteRow(db, folder_id, userId)

	// folder, err := repo.GetRow(db, folder_id, userId)

	// if err != nil {
	// 	fmt.Println("删除folder时取folder失败")
	// 	utils.Log(err)
	// 	return err
	// }

	images, err := imageCtrl.ItemsByFolder(db, folder_id, userId)

	if err != nil {
		fmt.Println("删除folder时取下属images失败")
		utils.Log(err)
		return err
	}

	// fmt.Println("images", images)

	for key := range images {
		err = imageCtrl.Delete(db, images[key].ID.Int, userId)
		if err != nil {
			fmt.Println("删除folder时删除下属images失败")
			utils.Log(err)
			return err
		}
	}

	return err
}
