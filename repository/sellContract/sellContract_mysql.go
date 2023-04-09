package sellContractRepository

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	buyContractRepository "github.com/xmluozp/creinox_server/repository/buyContract"
	commonRepo "github.com/xmluozp/creinox_server/repository/commonItem"
	financialTransactionRepository "github.com/xmluozp/creinox_server/repository/financialTransaction"
	orderFormRepo "github.com/xmluozp/creinox_server/repository/orderForm"
	portRepo "github.com/xmluozp/creinox_server/repository/port"
	sellSubitemRepository "github.com/xmluozp/creinox_server/repository/sellSubitem"
	userLogRepository "github.com/xmluozp/creinox_server/repository/userLog"

	"github.com/xmluozp/creinox_server/enums"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.SellContract
type repositoryName = Repository

var tableName = "sell_contract"

// 合同和order合体的view，显示用
var combineName = "combine_sell_contract"

// var viewName = "view_sell_contract"

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

	buyContractRepository := buyContractRepository.Repository{}
	sellSubitemRepository := sellSubitemRepository.Repository{}
	financialTransactionRepository := financialTransactionRepository.Repository{}

	for rows.Next() {

		item.ScanRowsView(rows)

		// 根据销售合同的ID（所有合同的起点），去搜索对应的工厂采购合同
		buyContract_list, _, _ := buyContractRepository.GetRows_fromSellContract(mydb, item.ID.Int, userId)
		item.BuyContractList = buyContract_list

		// 根据销售合同的ID 搜索子订单
		subitem_list, _, _ := sellSubitemRepository.GetRows_fromSellContract(mydb, item.ID.Int, userId)
		item.SellSubitem = subitem_list

		// 根据合同的ID 搜索转账记录
		trans_list, _, _ := financialTransactionRepository.GetRows_fromOrderForm(mydb, item.Order_form_id.Int, userId)
		item.FinancialTransactionList = trans_list

		// 根据合同的ID 搜索对应采购合同的还账记录
		order_form_ids := "0"
		for i := 0; i < len(buyContract_list); i++ {
			order_form_ids += "," + strconv.Itoa(buyContract_list[i].Order_form_id.Int)
		}
		trans_list_buyContract, _, _ := financialTransactionRepository.GetRows_fromOrderForms(mydb, order_form_ids, userId)
		item.FinancialTransactionList_buyContract = trans_list_buyContract

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
	orderitem.ContractType = nulls.NewInt(enums.ContractType.SellContract) // 销售合同type是1
	orderitem.Code = item.Code
	orderitem.InvoiceCode = item.InvoiceCode
	orderitem.Payable = nulls.NewFloat32(0)
	orderitem.PayablePaid = nulls.NewFloat32(0)
	orderitem.Receivable = item.TotalPrice
	orderitem.ReceivablePaid = item.PaidPrice
	orderitem.Seller_company_id = item.Seller_company_id
	orderitem.Buyer_company_id = item.Buyer_company_id
	orderitem.SellerAddress = item.SellerAddress
	orderitem.BuyerAddress = item.BuyerAddress
	orderitem.IsDone = item.IsDone
	orderitem.Order_memo = item.Order_memo

	orderFormRepo := orderFormRepo.Repository{}
	orderItem, errInsert := orderFormRepo.AddRow(mydb, orderitem, userId)

	if errInsert != nil {
		utils.Log(errInsert, "添加合同整体信息失败")
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
		utils.Log(errId)
		return item, errId
	}

	// 记录日志
	newItem, _ := b.GetRow(mydb, item.ID.Int, userId)

	err := b.ToUserLog(mydb, enums.LogActions["c"], newItem, userId)

	return item, err
}

func (b repositoryName) UpdateRow(mydb models.MyDb, item modelName, userId int) (int64, error) {

	item.UpdateUser_id = nulls.NewInt(userId)
	result, updatedRow, err := utils.DbQueryUpdate(mydb, tableName, combineName, item)

	// 从旧的item里面读取id
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
	// orderitem.Payable = nulls.NewFloat32(0)
	// orderitem.PayablePaid = nulls.NewFloat32(0)
	// orderitem.Receivable = item.TotalPrice
	// orderitem.ReceivablePaid = item.PaidPrice
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
	newItem, err := b.GetRow(mydb, item.ID.Int, userId)

	if err != nil {
		fmt.Println("记录日志获取新数据失败", err)
		return 0, err
	}

	err = b.ToUserLog(mydb, enums.LogActions["u"], newItem, userId)

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

	if err != nil {
		return nil, err
	}

	sellSubitemRepository := sellSubitemRepository.Repository{}

	//----------如果打印子列表，需要取出来
	subitem_list, _, err := sellSubitemRepository.GetRows_fromSellContract(mydb, id, userId)
	item.SellSubitem = subitem_list

	if err != nil {
		utils.Log(err, "GetPrintSource 出错1")
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	utils.ModifyDataSourceList(ds, "subitem_list", "ds_totalPrice",
		func(subitem map[string]interface{}) string {
			num1, ok1 := subitem["unitPrice"].(float64)
			num2, ok2 := subitem["amount"].(float64)
			if ok1 && ok2 {
				strNum := fmt.Sprintf("%.2f", num1*num2)
				return strNum
			}
			return "错误数据"
		})

	commonRepo := commonRepo.Repository{}
	portRepo := portRepo.Repository{}

	port1, err := portRepo.GetRow(mydb, item.Departure_port_id.Int, userId)
	port2, err := portRepo.GetRow(mydb, item.Destination_port_id.Int, userId)
	currency, err := commonRepo.GetRow(mydb, item.Currency_id.Int, userId)
	shippingType, err := commonRepo.GetRow(mydb, item.ShippingType_id.Int, userId)
	pricingTerm, err := commonRepo.GetRow(mydb, item.PricingTerm_id.Int, userId)

	if err != nil {
		utils.Log(err, "GetPrintSource 出错2")
		return ds, err
	}

	// 出发港，目标港，币种
	ds["ds_departure_port"] = port1.EName.String
	ds["ds_destination_port"] = port2.EName.String
	ds["ds_currency"] = currency.Ename.String
	ds["ds_shippingType"] = shippingType.Ename.String
	ds["ds_pricingTerm"] = pricingTerm.Ename.String

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

// SELECT SUM(unitPrice * amount) AS view_totalPrice
// FROM sell_subitem
// WHERE sell_contract_id = ??

// # SELECT a.*, b.view_totalPrice
// # FROM sell_contract a
// # LEFT JOIN (
// # 	SELECT sell_contract_id, SUM(unitPrice * amount) AS view_totalPrice
// # 	FROM sell_subitem
// # 	GROUP BY sell_contract_id
// # ) b ON a.id = b.sell_contract_id;

func (b repositoryName) ToUserLog(mydb models.MyDb, action string, item modelName, userId int) error {

	memo := fmt.Sprintf(`
		ID:			%d
		合同号:		%s
		总价:		%.2f
		交货期:		%s`,
		item.ID.Int, item.Code.String, item.TotalPrice.Float32, utils.FormatDate(item.DeliverAt.Time))

	// logBefore, _ := json.Marshal(before)
	// logAfter, _ := json.Marshal(after)

	var userLog models.UserLog
	userLog.Type = nulls.NewString(tableName)
	userLog.FunctionName = nulls.NewString(action)
	userLog.Memo = nulls.NewString(memo)
	// userLog.SnapshotBefore = nulls.NewString(string(logBefore))
	// userLog.SnapshotAfter = nulls.NewString(string(logAfter))

	_, err := userLogRepository.Repository{}.AddRow(mydb, userLog, userId)

	if err != nil {
		utils.Log(err, "日志出错")
		return err
	}

	return nil
}
