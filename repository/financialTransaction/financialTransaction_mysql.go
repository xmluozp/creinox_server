package financialTransactionRepository

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.FinancialTransaction
type repositoryName = Repository

var tableName = "financial_transaction"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// rows这里是一个cursor.
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

func (b repositoryName) GetRow(db *sql.DB, id int, userId int) (modelName, error) {
	var item modelName
	// row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)
	row := utils.DbQueryRow(db, "", tableName, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

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

	result, row, err := utils.DbQueryDelete(db, tableName, tableName, id, item)

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

func (b repositoryName) GetPrintSource(db *sql.DB, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(db, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
}
