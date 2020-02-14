package productRepository

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Product
type repositoryName = Repository

var tableName = "product"

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

	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	fmt.Println("add row", item)

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

func (b repositoryName) GetRows_DropDown(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	var sqlString string

	// 拦截 search.
	// isIncludeMeta: 包含meta
	isIncludeMeta := searchTerms["isIncludeMeta"]
	delete(searchTerms, "isIncludeMeta")

	if isIncludeMeta == "true" {
		// sqlString = fmt.Sprintf("SELECT * FROM %s WHERE path LIKE CONCAT((SELECT path FROM %s WHERE id = %d), ',' ,  %d , '%%')", tableName, tableName, root_id_int, root_id_int)
		sqlString = fmt.Sprintf("SELECT * FROM %s WHERE 1 = 1", tableName)

	} else {
		sqlString = fmt.Sprintf("SELECT a.* FROM %s a LEFT JOIN commodity_product b ON a.id = b.product_id WHERE b.isMeta <> 1", tableName)
	}

	rows, err := utils.DbQueryRows(db, sqlString, tableName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {
		err = rows.Scan(item.Receivers()...)
		if err != nil {
			fmt.Println("product scan出错", err.Error())
		}
		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}
