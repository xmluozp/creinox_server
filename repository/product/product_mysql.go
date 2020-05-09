package productRepository

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Product
type repositoryName = Repository

var tableName = "product"
var viewName = "view_product"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// ---customized:
	factory_id := searchTerms["companyFactory.id"]
	delete(searchTerms, "companyFactory.id")
	whereString := ""
	if factory_id != "" {
		whereString += fmt.Sprintf(" AND mainTable.id IN (SELECT product_id FROM product_purchase WHERE company_id = %s)", factory_id)
	}

	category_id := searchTerms["category_id"]
	delete(searchTerms, "category_id")

	if category_id != "" {
		whereString += fmt.Sprintf(" AND mainTable.category_id IN (SELECT id FROM category WHERE path LIKE '%%,%s' OR id = %s)", category_id, category_id)
	}

	rows, err := utils.DbQueryRows_Customized(db, "", tableName, &pagination, searchTerms, item, "", whereString)

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

	item.UpdateUser_id = nulls.NewInt(userId)
	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	fmt.Println("add row", item)

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
	result, err := utils.DbQueryUpdate(db, tableName, item)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (interface{}, error) {

	var item modelName

	result, row, err := utils.DbQueryDelete(db, tableName, id, item)

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

	return item, err
}

//---------------- customized
func (b repositoryName) GetRows_DropDown(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	var sqlString string

	// 拦截 search.
	// isIncludeMeta: 包含meta
	isIncludeMeta := searchTerms["isIncludeMeta"]
	delete(searchTerms, "isIncludeMeta")

	if isIncludeMeta == "1" {
		// sqlString = fmt.Sprintf("SELECT * FROM %s WHERE path LIKE CONCAT((SELECT path FROM %s WHERE id = %d), ',' ,  %d , '%%')", tableName, tableName, root_id_int, root_id_int)
		sqlString = fmt.Sprintf("SELECT mainTable.* FROM %s mainTable WHERE 1 = 1", viewName)

	} else {
		sqlString = fmt.Sprintf("SELECT mainTable.* FROM %s mainTable WHERE mainTable.id NOT IN (SELECT b.product_id FROM commodity_product b WHERE b.isMeta = 1)", viewName)
	}

	rows, err := utils.DbQueryRows(db, sqlString, viewName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {
		err = rows.Scan(item.ReceiversView()...)
		if err != nil {
			fmt.Println("product scan出错", err.Error())
		}
		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

// 根据销售合同取产品
func (b repositoryName) GetRows_DropDown_sellContract(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// 拦截 search.
	// isIncludeMeta: 包含meta
	sell_contract_id := searchTerms["sell_contract_id"]
	delete(searchTerms, "buy_contract_id")

	sqlString := fmt.Sprintf("SELECT mainTable.* FROM sell_subitem a LEFT JOIN commodity b ON a.commodity_id = b.id LEFT JOIN commodity_product c ON b.id = c.commodity_id LEFT JOIN %s mainTable ON c.product_id = mainTable.id WHERE a.sell_contract_id = %s",
		viewName,
		sell_contract_id)

	rows, err := utils.DbQueryRows(db, sqlString, viewName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	// ReceiversView
	for rows.Next() {
		fmt.Println(item.ReceiversView())
		err = rows.Scan(item.ReceiversView()...)
		if err != nil {
			fmt.Println("product scan出错", err.Error())
		}
		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRows_Component(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// 拦截 search.
	parent_id := searchTerms["parent_id"]
	delete(searchTerms, "parent_id")

	child_id := searchTerms["child_id"]
	delete(searchTerms, "child_id")

	subsql := "(SELECT * FROM " + tableName + ")"
	if parent_id != "" {
		subsql = fmt.Sprintf("(SELECT m1.* FROM "+tableName+" m1 LEFT JOIN product_component m2 ON m1.id = m2.child_id WHERE m2.parent_id = %s)", parent_id)
	} else if child_id != "" {
		subsql = fmt.Sprintf("(SELECT m1.* FROM "+tableName+" m1 LEFT JOIN product_component m2 ON m1.id = m2.parent_id WHERE m2.child_id = %s)", child_id)
	}

	rows, err := utils.DbQueryRows(db, "", subsql, &pagination, searchTerms, item)

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

func (b repositoryName) GetRows_ByCommodity(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	// 拦截 search.
	commodity_id := searchTerms["commodity_id"]
	delete(searchTerms, "commodity_id")

	// 不包括商品的元产品
	subsql := fmt.Sprintf("(SELECT m1.* FROM "+tableName+" m1 LEFT JOIN commodity_product m2 ON m1.id = m2.product_id WHERE m2.commodity_id = %s AND m2.isMeta = 0)", commodity_id)

	rows, err := utils.DbQueryRows(db, "", subsql, &pagination, searchTerms, item)

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

func (b repositoryName) Assemble(db *sql.DB, parent_id int, child_id int, userId int) error {

	_, err := db.Exec("INSERT INTO product_component (parent_id, child_id) VALUES(?, ?);", parent_id, child_id)

	return err

}

func (b repositoryName) Disassemble(db *sql.DB, parent_id int, child_id int, userId int) error {

	_, err := db.Exec("DELETE FROM product_component WHERE parent_id = ? AND child_id=?;", parent_id, child_id)

	return err
}
