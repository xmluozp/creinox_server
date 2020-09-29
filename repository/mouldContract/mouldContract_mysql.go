package mouldContractRepository

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/enums"
	"github.com/xmluozp/creinox_server/models"
	currencyRepo "github.com/xmluozp/creinox_server/repository/commonItem"
	financialTransactionRepository "github.com/xmluozp/creinox_server/repository/financialTransaction"
	orderFormRepo "github.com/xmluozp/creinox_server/repository/orderForm"
	productRepo "github.com/xmluozp/creinox_server/repository/product"
	userLogRepository "github.com/xmluozp/creinox_server/repository/userLog"

	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.MouldContract
type repositoryName = Repository

var tableName = "mould_contract"

// 合同和order合体的view，显示用
var combineName = "combine_mould_contract"

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

	image := models.Image{}
	financialTransactionRepository := financialTransactionRepository.Repository{}

	for rows.Next() {

		item.ScanRows(rows)

		item.View_image_thumbnail = nulls.NewString(image.AddPath(item.View_image_thumbnail.String))

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

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(mydb models.MyDb, item modelName, userId int) (modelName, error) {

	item.UpdateUser_id = nulls.NewInt(userId)

	// 抽出必要的字段，插入orderform，取出新生成的id
	orderitem := models.OrderForm{}
	orderitem.ContractType = nulls.NewInt(enums.ContractType.MouldContract) // 模板合同是3
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
	// orderresult, errInsert := utils.DbQueryInsert(mydb, tableName_order, orderitem)

	if errInsert != nil {
		return item, errInsert
	}

	// orderid, errId := orderresult.LastInsertId()
	item.Order_form_id = orderItem.ID
	// -------------------

	result, errInsert := utils.DbQueryInsert(mydb, tableName, item)

	if errInsert != nil {
		orderFormRepo.DeleteRow(mydb, orderItem.ID.Int, userId)
		utils.Log(nil, "添加合同详情失败，删除合同")
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
	orderitem.Payable = item.TotalPrice
	// orderitem.PayablePaid = item.PaidPrice
	orderitem.Receivable = nulls.NewFloat32(0)
	// orderitem.ReceivablePaid = nulls.NewFloat32(0)
	orderitem.Seller_company_id = item.Seller_company_id
	orderitem.Buyer_company_id = item.Buyer_company_id
	orderitem.SellerAddress = item.SellerAddress
	orderitem.BuyerAddress = item.BuyerAddress
	orderitem.IsDone = item.IsDone
	orderitem.Order_memo = item.Order_memo

	orderFormRepo := orderFormRepo.Repository{}
	_, err = orderFormRepo.UpdateRow(mydb, orderitem, userId)

	// result, row, err := utils.DbQueryUpdate(mydb, tableName_order, tableName_order, orderitem)
	// orderitem.ScanRow(row)

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
	// var orderitem models.OrderForm
	// result, row, err = utils.DbQueryDelete(mydb, tableName_order, tableName_order, item.Order_form_id.Int, orderitem)
	// orderitem.ScanRow(row)
	// -------
	orderFormRepo := orderFormRepo.Repository{}
	_, err = orderFormRepo.DeleteRow(mydb, item.Order_form_id.Int, userId)

	// 记录日志
	return item, err
}

func (b repositoryName) GetPrintSource(mydb models.MyDb, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	ds["ds_rmb"] = utils.FormatConvertNumToCny(item.TotalPrice.Float32)
	ds["activeAt"] = utils.FormatDate(item.ActiveAt.Time)
	ds["prepayAt"] = utils.FormatDate(item.PrepayAt.Time)
	ds["scheduleAt"] = utils.FormatDate(item.ScheduleAt.Time)
	ds["deliverAt"] = utils.FormatDate(item.DeliverAt.Time)
	ds["buyer_signAt"] = utils.FormatDate(item.Buyer_signAt.Time)
	ds["seller_signAt"] = utils.FormatDate(item.Seller_signAt.Time)

	// 取出一些比较深的数据
	productRepo := productRepo.Repository{}
	productRow, err := productRepo.GetRow(mydb, item.Product_id.Int, userId)

	if err != nil {
		return ds, err
	}

	currencyRepo := currencyRepo.Repository{}
	currencyRow, err := currencyRepo.GetRow(mydb, item.Currency_id.Int, userId)

	imageRow := productRow.ImageItem

	// 实验性的取数据。感觉还是orm好用

	// 约定格式： "path, width, height"
	ds["ds_image"] = imageRow.FileName.String + "," + strconv.Itoa(imageRow.Width.Int) + "," + strconv.Itoa(imageRow.Height.Int)
	ds["ds_product"] = productRow.Name.String
	ds["ds_currency"] = currencyRow.Name.String

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

func (b repositoryName) ToUserLog(mydb models.MyDb, action string, item modelName, userId int) {

	memo := fmt.Sprintf(`
		ID:			%d
		合同号:		%s
		总价:		%.2f
		签约日期:	 %s
		预付日期:	 %s
		预定交付期:	 %s
		实际交付期:	 %s`,
		item.ID.Int,
		item.Code.String,
		item.TotalPrice.Float32,
		utils.FormatDate(item.ActiveAt.Time),
		utils.FormatDate(item.PrepayAt.Time),
		utils.FormatDate(item.ScheduleAt.Time),
		utils.FormatDate(item.DeliverAt.Time))

	var userLog models.UserLog
	userLog.Type = nulls.NewString(tableName)
	userLog.FunctionName = nulls.NewString(action)
	userLog.Memo = nulls.NewString(memo)

	userLogRepository.Repository{}.AddRow(mydb, userLog, userId)
}
