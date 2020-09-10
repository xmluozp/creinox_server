package sellSubitemRepository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/enums"
	"github.com/xmluozp/creinox_server/models"
	orderFormRepo "github.com/xmluozp/creinox_server/repository/orderForm"
	userLogRepository "github.com/xmluozp/creinox_server/repository/userLog"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.SellSubitem
type repositoryName = Repository

var tableName = "sell_subitem"
var totalPriceName = "receivable" // 总价格是应收款还是应付款的总价格。设在代码开头，方便以后添加合同种类

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

	// 更新相应订单的总金额. 取出order_form_id
	order_form_id, err := b.getOrderFormId(db, item.ID.Int)
	if err != nil {
		return item, err
	}

	// 更新相应订单的总金额. 实际更新
	err = b.UpdateTotalPrice(db, order_form_id, userId)
	if err != nil {
		return item, err
	}

	// 记录日志
	var mapBefore map[string]interface{}
	mapAfter, _ := b.GetPrintSource(db, item.ID.Int, userId)
	b.ToUserLog(db, enums.LogActions["c"], mapBefore, mapAfter, item, userId)

	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	mapBefore, _ := b.GetPrintSource(db, item.ID.Int, userId)

	// 更新相应订单的总金额. 取出order_form_id
	order_form_id, err := b.getOrderFormId(db, item.ID.Int)
	if err != nil {
		return 0, err
	}

	// item.UpdateUser_id = nulls.NewInt(userId)
	result, row, err := utils.DbQueryUpdate(db, tableName, tableName, item)
	item.ScanRow(row)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	// 更新相应订单的总金额. 实际更新
	err = b.UpdateTotalPrice(db, order_form_id, userId)

	if err != nil {
		return 0, err
	}

	// 记录日志
	mapAfter, _ := b.GetPrintSource(db, item.ID.Int, userId)
	b.ToUserLog(db, enums.LogActions["u"], mapBefore, mapAfter, item, userId)

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (interface{}, error) {

	var item modelName

	// 记录日志
	mapBefore, _ := b.GetPrintSource(db, id, userId)
	var mapAfter map[string]interface{}
	b.ToUserLog(db, enums.LogActions["d"], mapBefore, mapAfter, item, userId)

	// 更新相应订单的总金额. 取出order_form_id
	order_form_id, err := b.getOrderFormId(db, id)
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

	// 更新相应订单的总金额. 实际更新
	err = b.UpdateTotalPrice(db, order_form_id, userId)
	if err != nil {
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

// =============================================== customized
// 每次item变动，都更新父合同里面的总价
// 根据item找到orderForm的id。随后更新总价格用
func (b repositoryName) getOrderFormId(db *sql.DB, id int) (order_form_id int, err error) {

	// 取出price和order form id
	row := db.QueryRow("SELECT a.order_form_id FROM sell_contract a LEFT JOIN sell_subitem b ON a.id = b.sell_contract_id WHERE b.id=?", id)
	err = row.Scan(&order_form_id)
	return order_form_id, err
}

// 更新特定order_form下面的销售合同总价格
func (b repositoryName) UpdateTotalPrice(db *sql.DB, order_form_id int, userId int) error {

	//tableName_order
	var totalPrice nulls.Float32

	// 取出price和order form id
	// row := db.QueryRow(`
	// SELECT b.view_totalPrice
	// FROM sell_contract a
	// LEFT JOIN (SELECT sell_contract_id, SUM(unitPrice * amount) AS view_totalPrice
	// 	FROM sell_subitem GROUP BY sell_contract_id) b ON a.id = b.sell_contract_id
	// LEFT JOIN sell_subitem c ON c.sell_contract_id = a.id  WHERE a.order_form_id=?`, order_form_id)

	row := db.QueryRow(
		`SELECT a.view_totalPrice FROM
			(SELECT sell_contract_id, SUM(unitPrice * amount) AS view_totalPrice FROM sell_subitem GROUP BY sell_contract_id) a 
			RIGHT JOIN sell_contract b
			ON a.sell_contract_id = b.id WHERE b.order_form_id = ?`, order_form_id)

	err := row.Scan(&totalPrice)

	if err != nil {
		return err
	}

	orderitem := models.OrderForm{}
	orderitem.ID = nulls.NewInt(order_form_id)
	orderitem.Receivable = nulls.NewFloat32(totalPrice.Float32) // 不convert一下，会提交null，然后被utils筛掉

	orderFormRepo := orderFormRepo.Repository{}
	_, err = orderFormRepo.UpdateRow(db, orderitem, userId)

	return err
}

// 根据合同号取出对应的item
func (b repositoryName) GetRows_fromSellContract(
	db *sql.DB,
	sell_contract_id int,
	userId int) ([]modelName, models.Pagination, error) {

	// var item modelName
	// var items []modelName
	var pagination models.Pagination
	searchTerms := make(map[string]string)

	// 不分页
	pagination.PerPage = -1

	sell_contract_id_str := strconv.Itoa(sell_contract_id)
	searchTerms["sell_contract_id"] = sell_contract_id_str

	// 这个应该是取出所有
	return b.GetRows(db, pagination, searchTerms, userId)
}

func (b repositoryName) ToUserLog(db *sql.DB, action string, before map[string]interface{}, after map[string]interface{}, item modelName, userId int) {

	memo := fmt.Sprintf(`
		ID:			%d
		产品:   	%s
		单价:		%.2f
		数量:		%d
		小计：		%.2f`,
		item.ID.Int,
		fmt.Sprintf(`[%s] %s`, item.Commodity.Code.String, item.Commodity.EName.String),
		item.UnitPrice.Float32,
		item.Amount.Int,
		item.UnitPrice.Float32*float32(item.Amount.Int))

	logBefore, _ := json.Marshal(before)
	logAfter, _ := json.Marshal(after)

	var userLog models.UserLog
	userLog.Type = nulls.NewString(tableName)
	userLog.FunctionName = nulls.NewString(action)
	userLog.Memo = nulls.NewString(memo)
	userLog.SnapshotBefore = nulls.NewString(string(logBefore))
	userLog.SnapshotAfter = nulls.NewString(string(logAfter))

	userLogRepository.Repository{}.AddRow(db, userLog, userId)
}
