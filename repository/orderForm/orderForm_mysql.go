package orderFormRepository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/enums"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"

	financialVoucherRepo "github.com/xmluozp/creinox_server/repository/financialVoucher"
)

type Repository struct{}
type modelName = models.OrderForm
type repositoryName = Repository

var tableName = "order_form"

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

	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		utils.Log(errInsert, "添加合同出错")
		return item, errInsert
	}

	id, err := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))

	if err != nil {
		utils.Log(err, "添加合同出错")
		return item, err
	}

	defer b.DeleteRow(db, item.ID.Int, userId)

	// =========================== 生成 voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	voucherItem1, voucherItem2 := b.getVoucher(db, item)
	_, err = financialVoucherRepo.AddRow(db, voucherItem1, userId)
	_, err = financialVoucherRepo.AddRow(db, voucherItem2, userId)

	return item, err
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	result, row, err := utils.DbQueryUpdate(db, tableName, tableName, item)
	item.ScanRow(row)

	if err != nil {
		utils.Log(err, "更新合同出错")
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		utils.Log(err, "更新合同出错")
		return 0, err
	}

	// =========================== 修改 voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	debit, credit := b.getVoucher(db, item)
	_, err = financialVoucherRepo.UpdateVoucher(db, debit, credit, userId)

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

	// =========================== 删除 voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	voucherItem1, voucherItem2 := b.getVoucher(db, item)
	_, err = financialVoucherRepo.DeleteVoucher(db, voucherItem1.Resource_code.String, userId)
	_, err = financialVoucherRepo.DeleteVoucher(db, voucherItem2.Resource_code.String, userId)

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

// ================= customized

func (b repositoryName) GetRows_DropDown(
	db *sql.DB,
	pagination models.Pagination, // 需要返回总页数
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

// 合同修改的时候，价格会变动，所以需要同步票据里的价格。以下根据resource code和item生成需要修改的票据里的value
func (b repositoryName) getVoucher(db *sql.DB, item modelName) (debit models.FinancialVoucher, credit models.FinancialVoucher) {

	// 返回的item1是借，item2是贷
	Debit_financialLedger_id := enums.FinancialLedgerType.UnDecided
	Credit_financialLedger_id := enums.FinancialLedgerType.UnDecided

	voucher_resource_code := fmt.Sprintf("order_form/%d", item.ID.Int)

	// *****合同类别区分***** 总账是应收还是应付?
	inout := ""
	switch item.ContractType.Int {
	case enums.ContractType.SellContract:
		inout = "in"
	case enums.ContractType.BuyContract:
		inout = "out"
	case enums.ContractType.MouldContract:
		inout = "out"
	}

	textMemo := ""
	switch inout {
	case "in":
		Debit_financialLedger_id = enums.FinancialLedgerType.ReceivableDebit
		Credit_financialLedger_id = enums.FinancialLedgerType.ReceivableCredit

		// debit.Debit = item.Receivable
		// credit.Credit = item.Receivable

		textMemo = "应收%s货款"
	case "out":
		Debit_financialLedger_id = enums.FinancialLedgerType.PayableDebit
		Credit_financialLedger_id = enums.FinancialLedgerType.PayableCredit

		// debit.Debit = item.Payable
		// credit.Credit = item.Payable

		textMemo = "应付%s货款"
	}

	amount := math.Max(float64(item.Payable.Float32), float64(item.Receivable.Float32))

	debit.Debit = nulls.NewFloat32(float32(amount))
	debit.Credit = nulls.NewFloat32(0)

	credit.Credit = nulls.NewFloat32(float32(amount))
	credit.Debit = nulls.NewFloat32(0)

	voucher_memo := fmt.Sprintf(textMemo, item.Code.String)

	debit.Resource_code = nulls.NewString(voucher_resource_code)
	debit.FinancialLedger_id = nulls.NewInt(Debit_financialLedger_id)
	debit.Memo = nulls.NewString(voucher_memo)

	credit.Resource_code = nulls.NewString(voucher_resource_code)
	credit.FinancialLedger_id = nulls.NewInt(Credit_financialLedger_id)
	credit.Memo = nulls.NewString(voucher_memo)

	return debit, credit
}
