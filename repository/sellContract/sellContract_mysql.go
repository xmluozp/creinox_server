package sellContractRepository

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	buyContractRepository "github.com/xmluozp/creinox_server/repository/buyContract"
	sellSubitemRepository "github.com/xmluozp/creinox_server/repository/sellSubitem"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.SellContract
type repositoryName = Repository

var tableName = "sell_contract"
var viewName = "view_sell_contract"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// rows这里是一个cursor.
	rows, err := utils.DbQueryRows(db, "", viewName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	buyContractRepository := buyContractRepository.Repository{}
	sellSubitemRepository := sellSubitemRepository.Repository{}

	for rows.Next() {

		item.ScanRowsView(rows)

		// 根据销售合同的ID（所有合同的起点），去搜索对应的工厂采购合同
		buyContract_list, _, _ := buyContractRepository.GetRows_fromSellContract(db, item.ID.Int, userId)
		item.BuyContractList = buyContract_list

		subitem_list, _, _ := sellSubitemRepository.GetRows_fromSellContract(db, item.ID.Int, userId)
		item.SellSubitem = subitem_list

		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int, userId int) (modelName, error) {
	var item modelName
	row := utils.DbQueryRow(db, "", viewName, id, item)

	err := item.ScanRowView(row)

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

// =============================================== customized

func (b repositoryName) GetRow_GetLast(db *sql.DB, id int, userId int) (modelName, error) {

	sqlstr := "SELECT * FROM " + tableName + " ORDER BY updateAt DESC LIMIT 1"

	var item modelName
	row := utils.DbQueryRow(db, sqlstr, viewName, 0, item)

	err := row.Scan(item.Receivers()...)

	return item, err
}
