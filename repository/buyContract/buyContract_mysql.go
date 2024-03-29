package buyContractRepository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/enums"
	"github.com/xmluozp/creinox_server/models"
	buySubitemRepository "github.com/xmluozp/creinox_server/repository/buySubitem"
	financialTransactionRepository "github.com/xmluozp/creinox_server/repository/financialTransaction"
	orderFormRepo "github.com/xmluozp/creinox_server/repository/orderForm"
	userLogRepository "github.com/xmluozp/creinox_server/repository/userLog"

	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.BuyContract
type repositoryName = Repository

var tableName = "buy_contract"

// 合同和order合体的view，显示用
var combineName = "combine_buy_contract"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	mydb models.MyDb,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// rows这里是一个cursor.
	rows, err := utils.DbQueryRows(mydb, "", combineName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	financialTransactionRepository := financialTransactionRepository.Repository{}

	for rows.Next() {

		item.ScanRowsView(rows)

		trans1_list, _, _ := financialTransactionRepository.GetRows_fromOrderForm(mydb, item.Order_form_id.Int, userId)
		item.FinancialTransactionList = trans1_list

		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(mydb models.MyDb, id int, userId int) (modelName, error) {
	var item modelName
	row := utils.DbQueryRow(mydb, "", combineName, id, item)

	err := item.ScanRowView(row)

	return item, err
}

func (b repositoryName) AddRow(mydb models.MyDb, item modelName, userId int) (modelName, error) {

	item.UpdateUser_id = nulls.NewInt(userId)

	// 抽出必要的字段，插入orderform，取出新生成的id
	orderitem := models.OrderForm{}
	orderitem.ContractType = nulls.NewInt(enums.ContractType.BuyContract) // 采购合同type是2
	orderitem.Code = item.Code
	orderitem.InvoiceCode = item.InvoiceCode
	orderitem.Payable = item.TotalPrice
	orderitem.PayablePaid = item.PaidPrice
	orderitem.Receivable = nulls.NewFloat32(0)
	orderitem.ReceivablePaid = nulls.NewFloat32(0)
	orderitem.Seller_company_id = item.Seller_company_id
	orderitem.Buyer_company_id = item.Buyer_company_id
	orderitem.SellerAddress = item.SellerAddress
	orderitem.BuyerAddress = item.BuyerAddress
	orderitem.IsDone = item.IsDone
	orderitem.Order_memo = item.Order_memo

	orderFormRepo := orderFormRepo.Repository{}
	orderItem, errInsert := orderFormRepo.AddRow(mydb, orderitem, userId)

	if errInsert != nil {
		return item, errInsert
	}

	item.Order_form_id = orderItem.ID

	// -------------------

	result, errInsert := utils.DbQueryInsert(mydb, tableName, item)

	if errInsert != nil {
		// orderFormRepo.DeleteRow(mydb, orderItem.ID.Int, userId)
		// utils.Log(nil, "添加合同详情失败，删除合同")
		utils.Log(errInsert, "添加合同详情失败")
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))
	if errId != nil {
		return item, errId
	}

	// 记录日志
	newItem, _ := b.GetRow(mydb, item.ID.Int, userId)
	b.ToUserLog(mydb, enums.LogActions["c"], newItem, userId)

	return item, errId
}

func (b repositoryName) UpdateRow(mydb models.MyDb, item modelName, userId int) (int64, error) {

	item.UpdateUser_id = nulls.NewInt(userId)

	result, updatedRow, err := utils.DbQueryUpdate(mydb, tableName, combineName, item)

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
	orderitem.Code = item.Code
	orderitem.InvoiceCode = item.InvoiceCode
	// orderitem.Payable = item.TotalPrice
	// orderitem.PayablePaid = item.PaidPrice
	// orderitem.Receivable = nulls.NewFloat32(0)
	// orderitem.ReceivablePaid = nulls.NewFloat32(0)
	orderitem.Seller_company_id = item.Seller_company_id
	orderitem.Buyer_company_id = item.Buyer_company_id
	orderitem.SellerAddress = item.SellerAddress
	orderitem.BuyerAddress = item.BuyerAddress
	orderitem.IsDone = item.IsDone
	orderitem.Order_memo = item.Order_memo

	orderFormRepo := orderFormRepo.Repository{}
	_, err = orderFormRepo.UpdateRow(mydb, orderitem, userId)

	if err != nil {
		return 0, err
	}
	// -------------------

	// 记录日志
	newItem, _ := b.GetRow(mydb, item.ID.Int, userId)
	b.ToUserLog(mydb, enums.LogActions["u"], newItem, userId)

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(mydb models.MyDb, id int, userId int) (interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	result, err := utils.DbQueryDelete(mydb, tableName, combineName, id, item)

	if err != nil {
		return nil, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil || rowsDeleted == 0 {
		return nil, err
	}

	// 删掉对应的order
	orderFormRepo := orderFormRepo.Repository{}
	_, err = orderFormRepo.DeleteRow(mydb, item.Order_form_id.Int, userId)

	// 记录日志
	b.ToUserLog(mydb, enums.LogActions["d"], item, userId)

	return item, err
}

func (b repositoryName) GetPrintSource(mydb models.MyDb, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	buySubitemRepository := buySubitemRepository.Repository{}

	//----------如果打印子列表，需要取出来
	subitem_list, _, err := buySubitemRepository.GetRows_fromBuyContract(mydb, id, userId)
	item.BuySubitem = subitem_list

	if err != nil {
		return nil, err
	}

	//----------在这里篡改需要打印的东西
	ds, err := utils.GetPrintSourceFromInterface(item)

	// 相乘获得总价格
	utils.ModifyDataSourceList(ds, "buy_subitem_list", "ds_totalPrice",
		func(subitem map[string]interface{}) string {
			num1, ok1 := subitem["unitPrice"].(float64)
			num2, ok2 := subitem["amount"].(float64)

			if ok1 && ok2 {
				strNum := fmt.Sprintf("%.2f", num1*num2)
				return strNum
			}
			return "错误数据"
		})

	// 提货时间的格式
	utils.ModifyDataSourceList(ds, "buy_subitem_list", "pickuptimeAt",
		func(subitem map[string]interface{}) string {

			pickupTime, ok := subitem["pickuptimeAt"].(string)

			if ok {
				t, err := time.Parse(time.RFC3339, pickupTime)

				if err != nil {
					return "错误数据"
				}

				return utils.FormatDate(t)
			}
			return ""
		})

	ds["ds_rmb"] = utils.FormatConvertNumToCny(item.TotalPrice.Float32)
	ds["activeAt"] = utils.FormatDate(item.ActiveAt.Time)

	return ds, err
}

// =============================================== customized

func (b repositoryName) GetRow_GetLast(mydb models.MyDb, id int, userId int) (modelName, error) {

	sqlstr := "SELECT * FROM " + combineName + " ORDER BY updateAt DESC LIMIT 1"

	var item modelName
	row := utils.DbQueryRow(mydb, sqlstr, combineName, 0, item)

	err := row.Scan(item.Receivers()...)

	return item, err
}

// 用来在销售合同界面，显示下属的采购合同列表
func (b repositoryName) GetRows_fromSellContract(
	mydb models.MyDb,
	sell_contract_id int,
	userId int) ([]modelName, models.Pagination, error) {

	// var item modelName
	// var items []modelName
	var pagination models.Pagination
	searchTerms := make(map[string]string)

	pagination.PerPage = -1
	// sell_contract_id_str := strconv.Itoa(sell_contract_id)
	sell_contract_id_str := strconv.Itoa(sell_contract_id)
	searchTerms["sell_contract_id"] = sell_contract_id_str

	// 这个应该是取出所有
	return b.GetRows(mydb, pagination, searchTerms, userId)
}

func (b repositoryName) ToUserLog(mydb models.MyDb, action string, item modelName, userId int) {

	memo := fmt.Sprintf(`
		ID:			%d
		合同号:		%s
		总价:		%.2f
		交货期:		%s`,
		item.ID.Int, item.Code.String, item.TotalPrice.Float32, utils.FormatDate(item.DeliverAt.Time))

	var userLog models.UserLog
	userLog.Type = nulls.NewString(tableName)
	userLog.FunctionName = nulls.NewString(action)
	userLog.Memo = nulls.NewString(memo)

	userLogRepository.Repository{}.AddRow(mydb, userLog, userId)
}
