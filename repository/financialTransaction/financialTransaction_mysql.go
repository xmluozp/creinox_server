package financialTransactionRepository

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/enums"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"

	financialAccountRepo "github.com/xmluozp/creinox_server/repository/financialAccount"
	financialVoucherRepo "github.com/xmluozp/creinox_server/repository/financialVoucher"

	orderFormRepo "github.com/xmluozp/creinox_server/repository/orderForm"
)

type Repository struct{}
type modelName = models.FinancialTransaction
type repositoryName = Repository

var tableName = "financial_transaction"

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

		item.ScanRows(rows)
		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(mydb models.MyDb, id int, userId int) (modelName, error) {
	var item modelName
	// row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)
	row := utils.DbQueryRow(mydb, "", tableName, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(mydb models.MyDb, item modelName, userId int) (modelName, error) {

	item.UpdateUser_id = nulls.NewInt(userId)
	result, errInsert := utils.DbQueryInsert(mydb, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))
	if errId != nil {
		return item, errId
	}

	// 只有添加才有side effect，修改转账记录禁止修改会导致side effect的字段
	err := b.sideEffects(mydb, item, userId)

	return item, err
}

func (b repositoryName) UpdateRow(mydb models.MyDb, item modelName, userId int) (int64, error) {

	// 部分列不允许修改，因为每条记录都有balance，如果金额、目标账户、目标合同改了，从今往后的balance就全错了：
	item.Amount_in = nulls.Float32{Float32: 0, Valid: false}
	item.Amount_out = nulls.Float32{Float32: 0, Valid: false}
	item.Balance = nulls.Float32{Float32: 0, Valid: false}
	item.Order_form_id = nulls.Int{Int: 0, Valid: false}
	item.FinancialAccount_id = nulls.Int{Int: 0, Valid: false}
	item.UpdateUser_id = nulls.NewInt(userId)

	result, row, err := utils.DbQueryUpdate(mydb, tableName, tableName, item)
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

func (b repositoryName) DeleteRow(mydb models.MyDb, id int, userId int) (interface{}, error) {

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

	// 去掉对应的金额。balance不处理。客户需要保证删除的是当前账户的最后一条记录
	err = b.sideEffectsReverse(mydb, item, userId)

	if err != nil {
		return nil, err
	}

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

func (b repositoryName) GetPrintSourceList(mydb models.MyDb, r *http.Request, userId int) (map[string]interface{}, error) {

	// 分页，排序，搜索关键词
	pagination := utils.GetPagination(r)
	searchTerms := utils.GetSearchTerms(r)

	items, _, err := b.GetRows(mydb, pagination, searchTerms, userId)

	// 统计数据
	totalIn := float32(0)
	totalOut := float32(0)

	for i := 0; i < len(items); i++ {
		totalIn += items[i].Amount_in.Float32
		totalOut += items[i].Amount_out.Float32
	}

	dataSource := make(map[string]interface{})
	dataSource["ds_list"] = items
	dataSource["ds_totalIn"] = totalIn
	dataSource["ds_totalOut"] = totalOut

	if len(items) > 0 {
		dataSource["ds_financialAccount"] = items[0].FinancialAccount.Name
		dataSource["ds_balance"] = items[0].FinancialAccount.Balance
		dataSource["ds_now"] = utils.FormatDateTime(time.Now())
	}

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(dataSource)
	utils.ModifyDataSourceList(ds, "ds_list", "transdateAt",
		func(subitem map[string]interface{}) string {

			t, err := time.Parse(time.RFC3339, subitem["transdateAt"].(string))

			if err != nil {
				return "错误数据"
			}

			return utils.FormatDateTime(t)
		})

	return ds, err
}

// ================= customized

// 收付款以及子合同的收付款
func (b repositoryName) GetRows_fromOrderForms(
	mydb models.MyDb,
	order_form_ids string,
	userId int) ([]modelName, models.Pagination, error) {

	// 先取出所有对应合同，逗号隔开

	var pagination models.Pagination
	searchTerms := make(map[string]string)

	// 不分页
	pagination.PerPage = -1
	order_form_id_str := order_form_ids
	searchTerms["order_form_id"] = order_form_id_str

	fmt.Println("合同：", order_form_id_str)

	// 这个应该是取出所有
	return b.GetRows(mydb, pagination, searchTerms, userId)
}

// 根据合同号取出对应的item
func (b repositoryName) GetRows_fromOrderForm(
	mydb models.MyDb,
	order_form_id int,
	userId int) ([]modelName, models.Pagination, error) {

	var pagination models.Pagination
	searchTerms := make(map[string]string)

	// 不分页
	pagination.PerPage = -1

	order_form_id_str := strconv.Itoa(order_form_id)
	searchTerms["order_form_id"] = order_form_id_str

	// 这个应该是取出所有
	return b.GetRows(mydb, pagination, searchTerms, userId)
}

func (b repositoryName) getVoucher(mydb models.MyDb, item modelName) (
	debit models.FinancialVoucher,
	credit models.FinancialVoucher) {

	voucher_resource_code := fmt.Sprintf("financial_transaction/%d", item.ID.Int)

	// 如果是合同款，明细就是合同号
	// if item.IsContractPayment.Bool {

	// 取出合同
	// orderForm := item.OrderForm

	// 200810: 移到前台了. 摘要在transaction添加的时候transUse
	// voucher_financialLedgeId := enums.FinancialSubjectType.UnDecided

	// *****合同类别区分***** 不同合同类别，付款所对应的科目
	// switch orderForm.ContractType.Int {
	// case enums.ContractType.SellContract:
	// 	voucher_financialLedgeId = enums.FinancialSubjectType.ReceivablePay
	// case enums.ContractType.BuyContract:
	// 	voucher_financialLedgeId = enums.FinancialSubjectType.PayablePay
	// case enums.ContractType.MouldContract:
	// 	voucher_financialLedgeId = enums.FinancialSubjectType.PayablePay
	// }
	// voucherItem.FinancialLedger_id = nulls.NewInt(voucher_financialLedgeId)

	// voucher_memo = fmt.Sprintf("%s %s", enums.ContractTypeLabel[orderForm.ContractType.Int], orderForm.Code.String)

	// } else { // 如果不是合同，摘要就是填写的内容

	// }

	debit.FinancialLedger_id = item.FinancialLedgerDebit_id
	credit.FinancialLedger_id = item.FinancialLedgerCredit_id

	amount := math.Max(float64(item.Amount_out.Float32), float64(item.Amount_in.Float32))

	// 借和贷都生成，根据 科目是否为空或者为0, 来决定是否添加
	debit.Debit = nulls.NewFloat32(float32(amount))
	debit.Credit = nulls.NewFloat32(0)

	credit.Credit = nulls.NewFloat32(float32(amount))
	credit.Debit = nulls.NewFloat32(0)

	debit.Memo = item.Tt_transUse
	credit.Memo = item.Tt_transUse

	debit.Resource_code = nulls.NewString(voucher_resource_code)
	credit.Resource_code = nulls.NewString(voucher_resource_code)

	debit.FinancialAccount_id = item.FinancialAccount_id
	credit.FinancialAccount_id = item.FinancialAccount_id

	return debit, credit
}

// 删除一条错误的记录
func (b repositoryName) sideEffectsReverse(mydb models.MyDb, item modelName, userId int) error {

	// ======================= 取消sideeffect
	// 如果是合同，已付款复原
	if item.IsContractPayment.Bool {

		orderFormRepo := orderFormRepo.Repository{}
		orderForm, err := orderFormRepo.GetRow(mydb, item.Order_form_id.Int, userId)
		if err != nil {
			fmt.Println("转账删除后连锁反应，取合同出错", err)
			return err
		}

		// 更新合同的已收付款
		orderForm.ReceivablePaid = nulls.NewFloat32(item.Amount_in.Float32 - orderForm.ReceivablePaid.Float32)
		orderForm.PayablePaid = nulls.NewFloat32(item.Amount_out.Float32 - orderForm.PayablePaid.Float32)

		_, err = orderFormRepo.UpdateRow(mydb, orderForm, userId)

		if err != nil {
			fmt.Println("更新合同出错", err)
			return err
		}
	}

	// ===================== 账号的balance还原
	financialAccountRepo := financialAccountRepo.Repository{}
	financialAccountItem, err := financialAccountRepo.GetRow(mydb, item.FinancialAccount_id.Int, userId)

	if err != nil {
		fmt.Println("连锁反应取账号合同出错", err)
		return err
	}
	oldBalance := financialAccountItem.Balance.Float32
	newBalance := oldBalance - item.Amount_in.Float32 + item.Amount_out.Float32
	financialAccountItem.Balance = nulls.NewFloat32(newBalance)
	_, err = financialAccountRepo.UpdateRow(mydb, financialAccountItem, userId)

	if err != nil {
		fmt.Println("更新transaction时，更新账户余额出错", err)
		return err
	}

	// ===================== 删掉 voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	voucherItemDebit, voucherItemCredit := b.getVoucher(mydb, item)

	_, err = financialVoucherRepo.DeleteVoucher(mydb, voucherItemDebit.Resource_code.String, userId)
	_, err = financialVoucherRepo.DeleteVoucher(mydb, voucherItemCredit.Resource_code.String, userId)
	return err
}

// 增加transaction之后触发的其他改动： 记录本身，3个balance：合同收付款，账户balance，transaction的balance.
// 如果是合同付款，强行指定目标公司。前台也不让修改
func (b repositoryName) sideEffects(mydb models.MyDb, item modelName, userId int) error {

	// ===================== 如果是针对合同付款的交易，更新合同本身以及生成对应的明细
	if item.IsContractPayment.Bool {

		orderFormRepo := orderFormRepo.Repository{}
		orderForm, err := orderFormRepo.GetRow(mydb, item.Order_form_id.Int, userId)
		if err != nil {
			fmt.Println("转账后连锁反应，取合同出错", err)
			return err
		}

		// *****合同类别区分***** 目标公司到底是合同里的甲方还是合同里的乙方，根据合同类型而定

		switch orderForm.ContractType.Int {
		case enums.ContractType.SellContract:
			item.Company_id = orderForm.Buyer_company_id

		case enums.ContractType.BuyContract:
			item.Company_id = orderForm.Seller_company_id

		case enums.ContractType.MouldContract:
			item.Company_id = orderForm.Seller_company_id
		}

		// 更新合同的已收付款
		orderForm.ReceivablePaid = nulls.NewFloat32(item.Amount_in.Float32 + orderForm.ReceivablePaid.Float32)
		orderForm.PayablePaid = nulls.NewFloat32(item.Amount_out.Float32 + orderForm.PayablePaid.Float32)

		_, err = orderFormRepo.UpdateRow(mydb, orderForm, userId)

		if err != nil {
			fmt.Println("更新合同出错", err)
			return err
		}
	}

	// ===================== 无论是不是合同，更新账户balance
	financialAccountRepo := financialAccountRepo.Repository{}
	financialAccountItem, err := financialAccountRepo.GetRow(mydb, item.FinancialAccount_id.Int, userId)

	if err != nil {
		fmt.Println("连锁反应取账号合同出错", err)
		return err
	}
	oldBalance := financialAccountItem.Balance.Float32
	newBalance := oldBalance + item.Amount_in.Float32 - item.Amount_out.Float32
	financialAccountItem.Balance = nulls.NewFloat32(newBalance)
	_, err = financialAccountRepo.UpdateRow(mydb, financialAccountItem, userId)
	if err != nil {
		fmt.Println("更新transaction时，更新账户余额出错", err)
		return err
	}

	// ===================== 通过update, 更新刚生成的transaction里的balance
	item.Balance = nulls.NewFloat32(newBalance)
	_, row, err := utils.DbQueryUpdate(mydb, tableName, tableName, item)
	item.ScanRow(row)

	if err != nil {
		fmt.Println("更新transaction的新余额出错", err)
		return err
	}

	// ===================== 生成 voucher。新建的时候，借贷科目如果没填，就不生成voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	voucherItemDebit, voucherItemCredit := b.getVoucher(mydb, item)

	if item.FinancialLedgerDebit_id.Valid {
		_, err = financialVoucherRepo.AddRow(mydb, voucherItemDebit, userId)
	}

	if item.FinancialLedgerCredit_id.Valid {
		_, err = financialVoucherRepo.AddRow(mydb, voucherItemCredit, userId)
	}

	return err

}
