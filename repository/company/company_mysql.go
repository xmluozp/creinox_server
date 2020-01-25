package companyRepoitory

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

	// 需要用join SELECT c.id, c.code, c.name, c.shortname, c.address, c.retrieveTime, u.userName FROM company c LEFT JOIN user u ON c.retriever_id = u.id
	rows, err := utils.DbQueryRows(db,
		"SELECT c.*, u.userName as 'retriever_id.userName' FROM "+
			tableName+" c LEFT JOIN user u ON c.retriever_id = u.id WHERE 1=1",
		tableName, &pagination, searchTerms, item)

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
	row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	// todo: 取图片
	fmt.Println("取公司", row, "SELECT * FROM "+tableName+" WHERE id = ?", id)
	// 假如不是平的struct而有子选项
	// 就要改写Scan
	// https://stackoverflow.com/questions/47335697/golang-decode-json-request-in-nested-struct-and-insert-in-db-as-blob

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

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

	item.ID = int(id)

	if errId != nil {
		return item, errId
	}

	// update folder
	folder.ID = int(folderId)
	folder.Memo = nulls.NewString("company/" + strconv.Itoa(item.ID))
	folder.FolderType = nulls.NewInt(1)
	folder.RefSource = nulls.NewString("company.gallary_folder_id")
	folder.RefId = nulls.NewInt(item.ID)

	result, errFolderUpdate := utils.DbQueryUpdate(db, "folder", folder)

	if errFolderUpdate != nil {
		return item, errFolderUpdate
	}

	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	item.UpdateUser_id = nulls.NewInt(userId)
	result, err := utils.DbQueryUpdate(db, tableName, item)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (int64, error) {

	// delete folder
	result, err := db.Exec("DELETE f, c FROM company c LEFT JOIN folder f ON f.id = c.gallary_folder_id WHERE c.id = ?", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, err
}
