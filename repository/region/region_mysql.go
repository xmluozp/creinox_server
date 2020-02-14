package regionRepository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Region
type repositoryName = Repository

var tableName = "region"

// =============================================== basic CRUD

func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination,
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	// 拦截 search
	root_id := searchTerms["root_id"]
	delete(searchTerms, "root_id")

	var sqlString string

	root_id_int, err := strconv.Atoi(root_id)

	// SELECT * FROM region WHERE path LIKE CONCAT((SELECT path FROM region WHERE id = 1), ',',1, '%') ORDER BY path ASC
	if err == nil && root_id_int > 0 {
		sqlString = fmt.Sprintf("SELECT * FROM %s WHERE path LIKE CONCAT((SELECT path FROM %s WHERE id = %d), ',' ,  %d , '%%')", tableName, tableName, root_id_int, root_id_int)
	} else {
		sqlString = ""
	}
	rows, err := utils.DbQueryRows(db, sqlString, tableName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {
		item.ScanRows(rows)
		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int) (modelName, error) {

	var item modelName
	row := utils.DbQueryRow(db, "", tableName, id, item)

	err := item.ScanRow(row)

	return item, err
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

	result, row, err := utils.DbQueryDelete(db, tableName, id, item)

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
