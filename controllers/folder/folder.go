package folderContrller

import (
	"database/sql"
	"fmt"

	imageController "github.com/xmluozp/creinox_server/controllers/imagedata"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/folder"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.Folder

var authName = ""

// =============================================== basic CRUD

// 删除folder把对应的image删掉
func (c Controller) Delete(db *sql.DB, id int, userId int) error {

	imageCtrl := imageController.Controller{}
	repo := repository.Repository{}

	folder, err := repo.GetRow(db, id)

	if err != nil {
		utils.Log(err)
		return err
	}

	images, err := imageCtrl.ItemsByFolder(db, folder.ID.Int)

	if err != nil {
		utils.Log(err)
		return err
	}

	fmt.Println("images", images)

	for key := range images {
		err = imageCtrl.Delete(db, images[key].ID.Int, userId)
		if err != nil {
			utils.Log(err)
			return err
		}
	}

	fmt.Println("images", db, id, userId)

	_, err = repo.DeleteRow(db, id, userId)

	return err
}
