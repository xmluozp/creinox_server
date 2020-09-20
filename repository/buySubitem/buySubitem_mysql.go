package buySubitemRepository

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
type modelName = models.BuySubitem
type repositoryName = Repository

var tableName = "buy_subitem"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	mydb models.MyDb,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// rows这里是一个cursor.
	rows, err := utils.DbQueryRows(mydb, "", tableName, &pagination, searchTerms, item)

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

func (b repositoryName) GetRow(mydb models.MyDb, id int, userId int) (modelName, error) {
	var item modelName
	row := utils.DbQueryRow(mydb, "", tableName, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(mydb models.MyDb, item modelName, userId int) (modelName, error) {

	// item.UpdateUser_id = nulls.NewInt(userId)
	result, errInsert := utils.DbQueryInsert(mydb, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))
	if errId != nil {
		return item, errId
	}

	// 更新相应订单的总金额. 取出order_form_id
	order_form_id, err := b.getOrderFormId(mydb, item.ID.Int)
	if err != nil {
		return item, err
	}

	// 更新相应订单的总金额
	err = b.UpdateTotalPrice(mydb, order_form_id, userId)
	if err != nil {
		return item, err
	}

	// 记录日志
	var mapBefore map[string]interface{}
	mapAfter, _ := b.GetPrintSource(mydb, item.ID.Int, userId)
	newItem, _ := b.GetRow(mydb, item.ID.Int, userId)
	b.ToUserLog(mydb, enums.LogActions["c"], mapBefore, mapAfter, newItem, userId)

	return item, errId
}

func (b repositoryName) UpdateRow(mydb models.MyDb, item modelName, userId int) (int64, error) {

	mapBefore, _ := b.GetPrintSource(mydb, item.ID.Int, userId)

	// 更新相应订单的总金额. 取出order_form_id
	order_form_id, err := b.getOrderFormId(mydb, item.ID.Int)
	if err != nil {
		return 0, err
	}

	result, row, err := utils.DbQueryUpdate(mydb, tableName, tableName, item)
	item.ScanRow(row)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	// 更新相应订单的总金额
	err = b.UpdateTotalPrice(mydb, order_form_id, userId)

	if err != nil {
		return 0, err
	}

	// 记录日志
	mapAfter, _ := b.GetPrintSource(mydb, item.ID.Int, userId)
	newItem, _ := b.GetRow(mydb, item.ID.Int, userId)
	b.ToUserLog(mydb, enums.LogActions["u"], mapBefore, mapAfter, newItem, userId)

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(mydb models.MyDb, id int, userId int) (interface{}, error) {

	mapBefore, _ := b.GetPrintSource(mydb, id, userId)

	// 更新相应订单的总金额. 取出order_form_id
	order_form_id, err := b.getOrderFormId(mydb, id)

	if err != nil {
		return nil, err
	}

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	result, err := utils.DbQueryDelete(mydb, tableName, tableName, id, item)

	if err != nil {
		return nil, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil || rowsDeleted == 0 {
		return nil, err
	}

	// 更新相应订单的总金额. 实际更新
	err = b.UpdateTotalPrice(mydb, order_form_id, userId)
	if err != nil {
		return nil, err
	}

	// 记录日志
	var mapAfter map[string]interface{}
	b.ToUserLog(mydb, enums.LogActions["d"], mapBefore, mapAfter, item, userId)

	return item, err
}

func (b repositoryName) GetPrintSource(mydb models.MyDb, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
}

// =============================================== customized

// 用做列表显示的collapse、以及打印的时候取子订单用
func (b repositoryName) GetRows_fromBuyContract(
	mydb models.MyDb,
	buy_contract_id int,
	userId int) ([]modelName, models.Pagination, error) {

	// var item modelName
	// var items []modelName
	var pagination models.Pagination
	searchTerms := make(map[string]string)

	// 不分页
	pagination.PerPage = -1

	buy_contract_id_str := strconv.Itoa(buy_contract_id)
	searchTerms["buy_contract_id"] = buy_contract_id_str

	returnitems, pagination, err := b.GetRows(mydb, pagination, searchTerms, userId)

	for i := 0; i < len(returnitems); i++ {
		returnitems[i].BuyContract = models.BuyContract{}
	}

	// 这个应该是取出所有
	return returnitems, pagination, err
}

// 每次item变动，都更新父合同里面的总价
// 根据item找到orderForm的id。随后更新总价格用
func (b repositoryName) getOrderFormId(mydb models.MyDb, id int) (order_form_id int, err error) {

	query := `
	SELECT a.order_form_id 
	FROM buy_contract a 
	LEFT JOIN buy_subitem b 
	ON a.id = b.buy_contract_id 
	WHERE b.id=?`

	var row *sql.Row
	// 取出price和order form id
	if mydb.Tx != nil {
		row = mydb.Tx.QueryRow(query, id)
	} else {
		row = mydb.Db.QueryRow(query, id)
	}

	err = row.Scan(&order_form_id)

	return order_form_id, err
}

// 每次item变动，都更新父合同里面的总价
func (b repositoryName) UpdateTotalPrice(mydb models.MyDb, order_form_id int, userId int) error {

	//tableName_order
	var totalPrice nulls.Float32

	// 取出price和order form id
	// row := db.QueryRow(`
	// SELECT a.order_form_id, b.view_totalPrice
	// FROM buy_contract a
	// LEFT JOIN (SELECT buy_contract_id, SUM(unitPrice * amount)
	// AS view_totalPrice FROM buy_subitem GROUP BY buy_contract_id) b ON a.id = b.buy_contract_id LEFT JOIN buy_subitem c ON c.buy_contract_id = a.id  WHERE c.id=?`, id)

	query := `SELECT a.view_totalPrice FROM
	(SELECT buy_contract_id, SUM(unitPrice * amount) AS view_totalPrice FROM buy_subitem GROUP BY buy_contract_id) a 
	RIGHT JOIN buy_contract b
	ON a.buy_contract_id = b.id WHERE b.order_form_id = ?`

	var row *sql.Row

	if mydb.Tx != nil {
		row = mydb.Tx.QueryRow(query, order_form_id)
	} else {
		row = mydb.Db.QueryRow(query, order_form_id)
	}

	err := row.Scan(&totalPrice)

	if err != nil {
		fmt.Println("错在更新价格", err)
		return err
	}

	orderitem := models.OrderForm{}
	orderitem.ID = nulls.NewInt(order_form_id)
	orderitem.Payable = nulls.NewFloat32(totalPrice.Float32) // 不convert一下，会提交null，然后被utils筛掉

	orderFormRepo := orderFormRepo.Repository{}
	_, err = orderFormRepo.UpdateRow(mydb, orderitem, userId)

	// _, err = db.Exec("UPDATE order_form SET "+totalPriceName+"=? WHERE id=?", &totalPrice, &order_form_id)
	return err
}

func (b repositoryName) ToUserLog(mydb models.MyDb, action string, before map[string]interface{}, after map[string]interface{}, item modelName, userId int) {

	memo := fmt.Sprintf(`
		ID:			%d
		产品:    	%s
		单价:		%.2f
		数量:		%d
		小计：		%.2f`,
		item.ID.Int,
		fmt.Sprintf(`[%s] %s`, item.Product.Code.String, item.Product.Name.String),
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

	userLogRepository.Repository{}.AddRow(mydb, userLog, userId)
}
