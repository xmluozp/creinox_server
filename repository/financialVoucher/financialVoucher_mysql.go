package financialVoucherRepository

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.FinancialVoucher
type repositoryName = Repository

var tableName = "financial_voucher"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

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
	row := utils.DbQueryRow(db, "", tableName, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	item.UpdateUser_id = nulls.NewInt(userId)

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

	item.UpdateUser_id = nulls.NewInt(userId)

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

// 根据voucher的code去升级，而不是根据id
func (b repositoryName) UpdateVoucher(db *sql.DB, item modelName, userId int) (int64, error) {

	// 根据code取出id
	sqlstr := "SELECT id FROM " + tableName + " WHERE resource_code = '" + item.Resource_code.String + "' LIMIT 1"
	row := utils.DbQueryRow(db, sqlstr, tableName, 0, item)

	var id int
	err := row.Scan(&id)

	if err != nil {
		// 如果没有对应的voucher，就忽略
		return 0, nil
	}

	// 修改传进来的item的id （原本是空的）
	item.ID = nulls.NewInt(id)

	return b.UpdateRow(db, item, userId)
}

// 根据voucher的code去删除，而不是根据id
func (b repositoryName) DeleteVoucher(db *sql.DB, voucherResourceCode string, userId int) (interface{}, error) {

	// 根据code取出id
	sqlstr := "SELECT id FROM " + tableName + " WHERE resource_code = '" + voucherResourceCode + "' LIMIT 1"
	var voucherItem modelName
	row := utils.DbQueryRow(db, sqlstr, tableName, 0, voucherItem)

	var id int
	err := row.Scan(&id)

	if err != nil {
		// 如果没有对应的voucher，就忽略
		return nil, nil
	}

	return b.DeleteRow(db, id, userId)
}
