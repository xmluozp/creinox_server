package imagedataRepository

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Image
type repositoryName = Repository

var tableName = "image"
var subsql = fmt.Sprintf("(SELECT m1.*, m2.memo FROM %s m1 LEFT JOIN folder m2 ON m1.gallary_folder_id = m2.id)", tableName)

// =============================================== basic CRUD

func (b repositoryName) GetRows(
	db *sql.DB,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 图片部分限制最多显示100张
	if pagination.PerPage < 0 {
		pagination.PerPage = 100
	}

	rows, err := utils.DbQueryRows(db, "", subsql, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {

		item.ScanRows(rows)
		items = append(items, item.Getter())
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int, userId int) (modelName, error) {

	var item modelName
	row := utils.DbQueryRow(db, "", subsql, id, item)

	err := item.ScanRow(row)
	return item.Getter(), err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	// result, errInsert := db.Exec("INSERT INTO role (name, rank, auth) VALUES(?, ?, ?);", item.Name, item.Rank, item.Auth)
	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))
	if errId != nil {
		return item, errId
	}

	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	result, row, err := utils.DbQueryUpdate(db, tableName, tableName, item)
	item.ScanRow(row)

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

	// customized
	rowDeleted := utils.DbQueryRow(db, "", subsql, id, item)
	// --- customized end

	var itemNotused modelName

	result, rowNotused, err := utils.DbQueryDelete(db, tableName, tableName, id, item)

	// 仅仅为了关闭连接
	itemNotused.ScanRow(rowNotused)

	if err != nil {
		return nil, err
	}

	err = item.ScanRow(rowDeleted)

	if err != nil {
		return nil, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil || rowsDeleted == 0 {
		return nil, err
	}

	return item, err
}

func (b repositoryName) GetPrintSource(db *sql.DB, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(db, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
}

func (b repositoryName) GetRowsByFolder(
	db *sql.DB,
	folderId int,
	userId int) (items []modelName, err error) {

	fmt.Println("folder?", folderId)

	// 需要用join SELECT a.runoob_id, a.runoob_author, b.runoob_count FROM runoob_tbl a INNER JOIN tcount_tbl b ON a.runoob_author = b.runoob_author;
	rows, err := db.Query("SELECT maintable.* FROM "+subsql+" maintable WHERE maintable.gallary_folder_id=?", folderId)

	if err != nil {
		return items, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {
		var item modelName
		item.ScanRows(rows)
		items = append(items, item)
	}

	return items, nil
}

// 因为要处理image，所以一个一个删
// DeleteRowMultiple
