package financialVoucherRepository

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.FinancialVoucher
type repositoryName = Repository

var tableName = "financial_voucher"

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
		utils.Log(errId)
		return item, errId
	}

	return item, errId
}

func (b repositoryName) UpdateRow(mydb models.MyDb, item modelName, userId int) (int64, error) {

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
	totalDebit := float32(0)
	totalCredit := float32(0)

	for i := 0; i < len(items); i++ {

		dsLedgeName := items[i].FinancialLedgerItem.LedgerName.String
		split := strings.Split(dsLedgeName, "/")

		if len(split) > 1 {
			items[i].FinancialLedgerItem.LedgerName = nulls.NewString(
				strings.Join(split[1:], "/"))
		}

		totalDebit += items[i].Debit.Float32
		totalCredit += items[i].Credit.Float32
	}

	dataSource := make(map[string]interface{})
	dataSource["ds_list"] = items
	dataSource["ds_totalDebit"] = totalDebit
	dataSource["ds_totalCredit"] = totalCredit

	if len(items) > 0 {
		dataSource["ds_now"] = utils.FormatDateTime(time.Now())
	}

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(dataSource)
	utils.ModifyDataSourceList(ds, "ds_list", "createAt",
		func(subitem map[string]interface{}) string {

			t, err := time.Parse(time.RFC3339, subitem["createAt"].(string))

			if err != nil {
				return "错误数据"
			}

			return utils.FormatDate(t)
		})

	// fmt.Println(ds)

	return ds, err
}

// 根据voucher的code去升级，而不是根据id
func (b repositoryName) UpdateVoucher(mydb models.MyDb, debit modelName, credit modelName, userId int) (rowsUpdated int64, err error) {

	// 根据code取出id。 会有两个，一借一贷。根据金额判断
	sqlstr := `SELECT id FROM %s WHERE financialLedger_id = %d AND resource_code = '%s' LIMIT 1`

	sqlCredit := fmt.Sprintf(sqlstr, tableName, credit.FinancialLedger_id.Int, credit.Resource_code.String)
	sqlDebit := fmt.Sprintf(sqlstr, tableName, debit.FinancialLedger_id.Int, debit.Resource_code.String)

	rowDebit := utils.DbQueryRow(mydb, sqlCredit, tableName, 0, debit)
	rowCredit := utils.DbQueryRow(mydb, sqlDebit, tableName, 0, credit)

	var idDebit int
	var idCredit int

	err = rowDebit.Scan(&idDebit)
	err = rowCredit.Scan(&idCredit)

	if err != nil {
		// 如果没有对应的voucher，就忽略
		return 0, nil
	}

	fmt.Println("更新凭证", debit, credit)

	// 修改传进来的item的id （原本是空的）
	debit.ID = nulls.NewInt(idDebit)
	credit.ID = nulls.NewInt(idCredit)

	rowsUpdated, err = b.UpdateRow(mydb, debit, userId)
	rowsUpdated, err = b.UpdateRow(mydb, credit, userId)

	return rowsUpdated, err
}

// 根据voucher的code去删除，而不是根据id
func (b repositoryName) DeleteVoucher(mydb models.MyDb, voucherResourceCode string, userId int) (interface{}, error) {

	// 根据code取出id
	sqlstr := "SELECT id FROM " + tableName + " WHERE resource_code = '" + voucherResourceCode + "' LIMIT 1"
	var voucherItem modelName
	row := utils.DbQueryRow(mydb, sqlstr, tableName, 0, voucherItem)

	var id int
	err := row.Scan(&id)

	if err != nil {
		// 如果没有对应的voucher，就忽略
		return nil, nil
	}

	return b.DeleteRow(mydb, id, userId)
}
