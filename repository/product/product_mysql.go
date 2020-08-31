package productRepository

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"

	commodityRepo "github.com/xmluozp/creinox_server/repository/commodity"
)

type Repository struct{}
type modelName = models.Product
type repositoryName = Repository

var tableName = "product"
var viewName = "view_product"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// ---customized:
	factory_id := searchTerms["companyFactory.id"]

	// 20200607: 不知道为啥要删掉，删掉了se就无法返回前端了。前端下一页就会出问题。目前方案在前端生成一个新的searchTerm用作返回
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
	result, row, err := utils.DbQueryUpdate(db, tableName, tableName, item)
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

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (interface{}, error) {

	var item modelName

	// 在删除前取出meta商品id备用 (因为删除后的cascade会导致关联信息消失)
	commodityRepo := commodityRepo.Repository{}
	commodityItem, err := commodityRepo.GetRow_ByProduct(db, id, userId)

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

	// 如果删除产品成功，也删除对应的元商品
	commodityRepo.DeleteRow(db, commodityItem.ID.Int, userId)

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

//---------------- customized
func (b repositoryName) GetRows_DropDown(
	db *sql.DB,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

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
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 search.
	sell_contract_id := searchTerms["sell_contract_id"]
	delete(searchTerms, "sell_contract_id")

	if sell_contract_id == "" { // 如果为空，就让它搜不到
		sell_contract_id = "-1"
	}

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

// 根据合同的子合同，去搜索子合同对应的商品，然后关联到下属产品
func (b repositoryName) GetRows_DropDown_sellSubitem(
	db *sql.DB,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	sell_subitem_id := searchTerms["sell_subitem_id"]
	delete(searchTerms, "sell_subitem_id")

	if sell_subitem_id == "" { // 如果为空，就让它搜不到
		sell_subitem_id = "-1"
	}

	sqlString := fmt.Sprintf("SELECT mainTable.* FROM sell_subitem a LEFT JOIN commodity b ON a.commodity_id = b.id LEFT JOIN commodity_product c ON b.id = c.commodity_id LEFT JOIN %s mainTable ON c.product_id = mainTable.id WHERE a.id = %s",
		viewName,
		sell_subitem_id)
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
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 search.
	parent_id := searchTerms["parent_id"]
	// delete(searchTerms, "parent_id")

	child_id := searchTerms["child_id"]
	// delete(searchTerms, "child_id")

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
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 search.
	commodity_id := searchTerms["commodity_id"]
	delete(searchTerms, "commodity_id")

	// 包括商品的元产品（避免用户迷惑，但还是需要标注出来）
	subsql := fmt.Sprintf("(SELECT m1.* FROM "+tableName+" m1 LEFT JOIN commodity_product m2 ON m1.id = m2.product_id WHERE m2.commodity_id = %s)", commodity_id)

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
