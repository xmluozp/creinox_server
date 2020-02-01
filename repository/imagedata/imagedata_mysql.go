package imagedataRepository

import (
	"database/sql"
	"fmt"

	"github.com/Unknwon/goconfig"
	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Image
type repositoryName = Repository

var tableName = "image"
var UPLOAD_FOLDER = "uploads/"

// =============================================== basic CRUD

func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	// 需要用join SELECT a.runoob_id, a.runoob_author, b.runoob_count FROM runoob_tbl a INNER JOIN tcount_tbl b ON a.runoob_author = b.runoob_author;
	rows, err := utils.DbQueryRows(db, "SELECT * FROM "+tableName+" WHERE 1=1 ", tableName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}

	rootUrl, err := cfg.GetValue("site", "root")
	port, err := cfg.Int("site", "port")
	uploadFolder := fmt.Sprintf("%s:%d/%s", rootUrl, port, UPLOAD_FOLDER)

	for rows.Next() {

		// 把数据库读出来的列填进对应的变量里 (如果只想取对应的列怎么办？)
		// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错
		item.ScanRows(rows)
		item.ThumbnailPath = nulls.NewString(UPLOAD_FOLDER + item.ThumbnailPath.String)
		item.Path = nulls.NewString(uploadFolder + item.Path.String)

		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	// for i, _ := range items {
	// 	items[i].Name = "改头换面"
	// }

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int) (modelName, error) {

	var item modelName
	row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	err := item.ScanRow(row)

	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}

	rootUrl, err := cfg.GetValue("site", "root")
	port, err := cfg.Int("site", "port")
	uploadFolder := fmt.Sprintf("%s:%d/%s", rootUrl, port, UPLOAD_FOLDER)

	item.ThumbnailPath = nulls.NewString(UPLOAD_FOLDER + item.ThumbnailPath.String)
	item.Path = nulls.NewString(uploadFolder + item.Path.String)

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	// result, errInsert := db.Exec("INSERT INTO role (name, rank, auth) VALUES(?, ?, ?);", item.Name, item.Rank, item.Auth)
	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = int(id)
	if errId != nil {
		return item, errId
	}

	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

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

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (interface{}, error) {

	var item modelName

	result, row, err := utils.DbQueryDelete(db, tableName, id)

	if err != nil {
		return nil, err
	}

	err = item.ScanRow(row)

	if err != nil {
		return nil, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil || rowsDeleted == 0 {
		return nil, err
	}

	return item, err
}

func (b repositoryName) GetRowsByFolder(
	db *sql.DB,
	folderId int) (items []modelName, err error) {

	// 需要用join SELECT a.runoob_id, a.runoob_author, b.runoob_count FROM runoob_tbl a INNER JOIN tcount_tbl b ON a.runoob_author = b.runoob_author;
	rows, err := db.Query("SELECT * FROM "+tableName+" WHERE gallary_folder_id=?", folderId)

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
