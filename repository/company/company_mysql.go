package companyRepository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Company
type repositoryName = Repository

var tableName = "company"

// =============================================== basic CRUD

func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	rows, err := utils.DbQueryRows(db, "", tableName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {

		err := item.ScanRows(rows)

		if err != nil {
			fmt.Println("scan err", err.Error())
		}
		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int) (modelName, error) {

	var item modelName
	// row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	// todo: 取图片
	row := utils.DbQueryRow(db, "", tableName, id, item)
	err := item.ScanRow(row)

	// imageCtrl := imageController.Controller{}

	// if item.ImageLicense_id.Valid {
	// 	license, err := imageCtrl.Item(db, item.ImageLicense_id.Int)
	// 	if err != nil {
	// 		return item, err
	// 	}
	// 	item.ImageLicense = license
	// }

	// if item.ImageBizCard_id.Valid {
	// 	bizcard, err := imageCtrl.Item(db, item.ImageBizCard_id.Int)
	// 	if err != nil {
	// 		return item, err
	// 	}
	// 	item.ImageBizCard = bizcard
	// }

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, itemRec interface{}, userId int) (modelName, error) {

	item := itemRec.(modelName)

	folder := models.Folder{}

	// company folder
	newFolder, errInsert := utils.DbQueryInsert(db, "folder", folder)
	folderId, err := newFolder.LastInsertId()

	if err != nil {
		return item, err
	}

	item.UpdateUser_id = nulls.NewInt(userId)
	item.Gallary_folder_id = nulls.NewInt(int(folderId))

	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()

	item.ID = nulls.NewInt(int(id))

	if errId != nil {
		return item, errId
	}

	// update folder
	folder.ID = nulls.NewInt(int(folderId))
	folder.Memo = nulls.NewString("company/" + strconv.Itoa(item.ID.Int))
	folder.FolderType = nulls.NewInt(1)
	folder.RefSource = nulls.NewString("company.gallary_folder_id")
	folder.RefId = item.ID

	result, errFolderUpdate := utils.DbQueryUpdate(db, "folder", folder)

	if errFolderUpdate != nil {
		return item, errFolderUpdate
	}

	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, itemRec interface{}, userId int) (int64, error) {

	fmt.Println("h", itemRec)
	item := itemRec.(modelName)

	item.UpdateUser_id = nulls.NewInt(userId)
	result, err := utils.DbQueryUpdate(db, tableName, item)

	// fmt.Println("database updateRow", item)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (interface{}, error) {

	var item modelName
	result, row, err := utils.DbQueryDelete(db, tableName, id, item)
	err = item.ScanRow(row)

	fmt.Println("scaned 扫描以后", item)
	// result, err := db.Exec("DELETE f, c FROM company c LEFT JOIN folder f ON f.id = c.gallary_folder_id WHERE c.id = ?", id)

	if err != nil {
		return nil, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil || rowsDeleted == 0 {
		return nil, err
	}

	return item, err
}