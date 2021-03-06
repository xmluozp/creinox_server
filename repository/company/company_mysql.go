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
	mydb models.MyDb,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	rows, err := utils.DbQueryRows(mydb, "", tableName, &pagination, searchTerms, item)

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

func (b repositoryName) GetRow(mydb models.MyDb, id int, userId int) (modelName, error) {

	var item modelName
	// row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	// todo: 取图片
	row := utils.DbQueryRow(mydb, "", tableName, id, item)
	err := item.ScanRow(row)

	// imageCtrl := imageController.Controller{}

	// if item.ImageLicense_id.Valid {
	// 	license, err := imageCtrl.Item(mydb, item.ImageLicense_id.Int)
	// 	if err != nil {
	// 		return item, err
	// 	}
	// 	item.ImageLicense = license
	// }

	// if item.ImageBizCard_id.Valid {
	// 	bizcard, err := imageCtrl.Item(mydb, item.ImageBizCard_id.Int)
	// 	if err != nil {
	// 		return item, err
	// 	}
	// 	item.ImageBizCard = bizcard
	// }

	return item, err
}

func (b repositoryName) AddRow(mydb models.MyDb, itemRec interface{}, userId int) (modelName, error) {

	item := itemRec.(modelName)

	folder := models.Folder{}

	// company folder
	newFolder, errInsert := utils.DbQueryInsert(mydb, "folder", folder)
	folderId, err := newFolder.LastInsertId()

	if err != nil {
		return item, err
	}

	item.UpdateUser_id = nulls.NewInt(userId)
	item.Gallary_folder_id = nulls.NewInt(int(folderId))

	result, errInsert := utils.DbQueryInsert(mydb, tableName, item)

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

	result, row, errFolderUpdate := utils.DbQueryUpdate(mydb, "folder", "folder", folder)
	folder.ScanRow(row)

	if errFolderUpdate != nil {
		return item, errFolderUpdate
	}

	return item, errId
}

func (b repositoryName) UpdateRow(mydb models.MyDb, itemRec interface{}, userId int) (int64, error) {

	fmt.Println("h", itemRec)
	item := itemRec.(modelName)

	item.UpdateUser_id = nulls.NewInt(userId)
	result, row, err := utils.DbQueryUpdate(mydb, tableName, tableName, item)
	item.ScanRow(row)

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

func (b repositoryName) DeleteRow(mydb models.MyDb, id int, userId int) (interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	result, err := utils.DbQueryDelete(mydb, tableName, tableName, id, item)

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

func (b repositoryName) GetPrintSource(mydb models.MyDb, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
}

// 公司专用，取当前类别，当前前缀的code里最大的那个返回
func (b repositoryName) GetRow_byCode(mydb models.MyDb, companyType int, keyWord string, userId int) (modelName, error) {

	var item modelName
	// row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	query := `SELECT company.*
	FROM company 
	WHERE UPPER(code) LIKE CONCAT(UPPER(?), "%") AND companyType = ?
	AND CONVERT(SUBSTRING(code, ?, length(code)), UNSIGNED) > 0
	ORDER BY CONVERT(SUBSTRING(code, ?, length(code)), UNSIGNED) DESC LIMIT 1`

	var row *sql.Row

	if mydb.Tx != nil {
		row = mydb.Tx.QueryRow(query, keyWord, companyType, len(keyWord)+1, len(keyWord)+1)
	} else {
		row = mydb.Db.QueryRow(query, keyWord, companyType, len(keyWord)+1, len(keyWord)+1)
	}

	err := row.Scan(item.Receivers()...)

	if err != nil {
		return item, err
	}

	return item, err
}
