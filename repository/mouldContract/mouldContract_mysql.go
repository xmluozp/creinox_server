package buyContractRepository

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.MouldContract
type repositoryName = Repository

var tableName = "mould_contract"

// 合同和order合体的view，显示用
var combineName = "combine_mould_contract"
var tableName_order = "order_form"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// rows这里是一个cursor.
	rows, err := utils.DbQueryRows(db, "", combineName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	image := models.Image{}

	for rows.Next() {

		item.ScanRows(rows)

		item.View_image_thumbnail = nulls.NewString(image.AddPath(item.View_image_thumbnail.String))

		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int, userId int) (modelName, error) {
	var item modelName
	row := utils.DbQueryRow(db, "", combineName, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	item.UpdateUser_id = nulls.NewInt(userId)

	// 抽出必要的字段，插入orderform，取出新生成的id
	orderitem := models.OrderForm{}
	orderitem.Type = nulls.NewInt(int(2)) // 采购合同type是2
	orderitem.TotalPrice = item.TotalPrice
	orderitem.PaidPrice = item.PaidPrice
	orderitem.Seller_company_id = item.Seller_company_id
	orderitem.Buyer_company_id = item.Buyer_company_id
	orderitem.IsDone = item.IsDone
	orderitem.Order_memo = item.Order_memo

	orderresult, errInsert := utils.DbQueryInsert(db, tableName_order, orderitem)

	if errInsert != nil {
		return item, errInsert
	}
	orderid, errId := orderresult.LastInsertId()
	item.Order_form_id = nulls.NewInt(int(orderid))
	// -------------------

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

	result, updatedRow, err := utils.DbQueryUpdate(db, tableName, combineName, item)

	var olditem modelName
	olditem.ScanRow(updatedRow)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	// 升级完顺便升级orderform
	orderitem := models.OrderForm{}
	orderitem.ID = olditem.Order_form_id
	orderitem.TotalPrice = item.TotalPrice
	orderitem.PaidPrice = item.PaidPrice
	orderitem.Seller_company_id = item.Seller_company_id
	orderitem.Buyer_company_id = item.Buyer_company_id
	orderitem.IsDone = item.IsDone
	orderitem.Order_memo = item.Order_memo

	result, _, err = utils.DbQueryUpdate(db, tableName_order, tableName_order, orderitem)

	if err != nil {
		return 0, err
	}
	// -------------------

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (interface{}, error) {

	var item modelName

	result, row, err := utils.DbQueryDelete(db, tableName, combineName, id, item)

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

	// 删掉对应的order
	var orderitem models.OrderForm
	result, row, err = utils.DbQueryDelete(db, tableName_order, tableName_order, item.Order_form_id.Int, orderitem)
	// -------

	return item, err
}

func (b repositoryName) GetPrintSource(db *sql.DB, id int, userId int) (modelName, error) {

	return b.GetRow(db, id, userId)
}

// =============================================== customized

func (b repositoryName) GetRow_GetLast(db *sql.DB, id int, userId int) (modelName, error) {

	sqlstr := "SELECT * FROM " + combineName + " ORDER BY updateAt DESC LIMIT 1"

	var item modelName
	row := utils.DbQueryRow(db, sqlstr, combineName, 0, item)

	err := row.Scan(item.Receivers()...)

	return item, err
}
