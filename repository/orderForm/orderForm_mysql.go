package orderFormRepository

import (
	"database/sql"
	"fmt"

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
		return item, errInsert
	}

	id, err := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))
	if err != nil {
		return item, err
	}

	// =========================== 生成 voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	voucherItem := b.getVoucher(db, item)
	_, err = financialVoucherRepo.AddRow(db, voucherItem, userId)

	if err != nil {
		return item, err
	}

	return item, err
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	result, row, err := utils.DbQueryUpdate(db, tableName, tableName, item)
	item.ScanRow(row)

	if err != nil {
		fmt.Println("错在更新row", err)
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	// =========================== 生成 voucher
	financialVoucherRepo := financialVoucherRepo.Repository{}
	voucherItem := b.getVoucher(db, item)
	_, err = financialVoucherRepo.UpdateVoucher(db, voucherItem, userId)

	if err != nil {
		fmt.Println("错在更新voucher", err)
		return 0, err
	}

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
	voucherItem := b.getVoucher(db, item)
	_, err = financialVoucherRepo.DeleteVoucher(db, voucherItem.Resource_code.String, userId)

	if err != nil {
		return 0, err
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

// ================= customized

func (b repositoryName) getVoucher(db *sql.DB, item modelName) models.FinancialVoucher {

	voucherItem := models.FinancialVoucher{}
	voucher_financialSubject := ""
	voucher_resource_code := fmt.Sprintf("order_form/%d", item.ID.Int)

	// 看总账是应收应付. 应付增加是贷，付款是借（在付款的地方写）
	switch item.ContractType.Int {
	case enums.ContractType.SellContract:
		voucher_financialSubject = enums.FinancialSubjectType.Receivable
		voucherItem.Debit = item.Receivable
	case enums.ContractType.BuyContract:
		voucher_financialSubject = enums.FinancialSubjectType.Payable
		voucherItem.Credit = item.Payable
	case enums.ContractType.MouldContract:
		voucher_financialSubject = enums.FinancialSubjectType.Payable
		voucherItem.Credit = item.Payable
	}

	voucher_detailSubject := fmt.Sprintf("%s %s", enums.ContractTypeLabel[item.ContractType.Int], item.Code.String)

	voucherItem.Memo = item.Order_memo
	voucherItem.FinancialSubject = nulls.NewString(voucher_financialSubject)
	voucherItem.Resource_code = nulls.NewString(voucher_resource_code)
	voucherItem.DetailedSubject = nulls.NewString(voucher_detailSubject)

	return voucherItem

}
