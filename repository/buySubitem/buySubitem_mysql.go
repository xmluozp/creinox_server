package buySubitemRepository

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.BuySubitem
type repositoryName = Repository

var tableName = "buy_subitem"

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

		err = item.ScanRows(rows)
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

	// item.UpdateUser_id = nulls.NewInt(userId)
	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()

	item.ID = nulls.NewInt(int(id))

	if errId != nil {
		return item, errId
	}

	// 更新相应订单的总金额
	err := b.UpdateTotalPrice(db, item.ID.Int, userId)
	if err != nil {
		return item, err
	}
	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	// 更新相应订单的总金额
	err := b.UpdateTotalPrice(db, item.ID.Int, userId)

	if err != nil {
		return 0, err
	}

	result, _, err := utils.DbQueryUpdate(db, tableName, tableName, item)

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

	// 更新相应订单的总金额(放前面因为删除了就没有了)
	err := b.UpdateTotalPrice(db, id, userId)

	if err != nil {
		return nil, err
	}

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

// =============================================== customized
// 每次item变动，都更新父合同里面的总价
func (b repositoryName) UpdateTotalPrice(db *sql.DB, id int, userId int) error {

	//tableName_order
	var totalPrice float32
	var order_form_id int

	// 取出price和order form id
	row := db.QueryRow("SELECT a.order_form_id, b.view_totalPrice FROM buy_contract a LEFT JOIN (SELECT buy_contract_id, SUM(unitPrice * amount) AS view_totalPrice FROM buy_subitem GROUP BY buy_contract_id) b ON a.id = b.buy_contract_id LEFT JOIN buy_subitem c ON c.buy_contract_id = a.id  WHERE c.id=?", id)

	err := row.Scan(&order_form_id, &totalPrice)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE order_form SET totalPrice=? WHERE id=?", &totalPrice, &order_form_id)
	return err
}
