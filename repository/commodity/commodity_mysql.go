package commodityRepository

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Commodity
type repositoryName = Repository

var tableName = "commodity"

// var tableMeta = "(SELECT core1.*, core2.product_id as product_id, core3.image_id as image_id FROM commodity core1 LEFT JOIN commodity_product core2 ON core1.id = core2.commodity_id LEFT JOIN product core3 ON core3.id = core2.product_id  WHERE core2.isMeta = 1)"

// 不判断meta，不然删了产品，导致关联删掉了以后，就成了幽灵记录
var tableMeta = "(SELECT core1.*, core2.product_id as product_id, core3.image_id as image_id FROM commodity core1 LEFT JOIN commodity_product core2 ON core1.id = core2.commodity_id LEFT JOIN product core3 ON core3.id = core2.product_id)"

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	// rows这里是一个cursor.
	rows, err := utils.DbQueryRows(db, "", tableMeta, &pagination, searchTerms, item)

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

func (b repositoryName) GetRow(db *sql.DB, id int) (modelName, error) {
	var item modelName

	row := utils.DbQueryRow(db, "", tableMeta, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) GetRow_ByProduct(db *sql.DB, id int) (modelName, error) {
	var item modelName

	rowCommodityProduct := db.QueryRow("SELECT commodity_id FROM commodity_product WHERE product_id = ? AND isMeta = 1", id)

	err := rowCommodityProduct.Scan(&item.ID)

	if err != nil {
		return item, err
	}

	row := utils.DbQueryRow(db, "", tableMeta, item.ID.Int, item)

	err = item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	// result, errInsert := db.Exec("INSERT INTO role (name, rank, auth) VALUES(?, ?, ?);", item.Name, item.Rank, item.Auth)
	item.UpdateUser_id = nulls.NewInt(userId)
	result, errInsert := utils.DbQueryInsert(db, tableName, item)

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

	// result, row, err := utils.DbQueryDelete(db, tableName, id, item)
	// 这里特殊，返回一个空的item。因为commodity的image_id字段是在关联的product表里，这里delete的地方取不到，而且外部用不到这个item
	result, _, err := utils.DbQueryDelete(db, tableName, id, item)

	if err != nil {
		return nil, err
	}

	// err = item.ScanRow(row)

	// if err != nil {
	// 	return nil, err
	// }

	rowsDeleted, err := result.RowsAffected()

	if err != nil || rowsDeleted == 0 {
		return nil, err
	}

	return item, err
}

//==================== customized
func (b repositoryName) AddRow_WithProduct(db *sql.DB, commodity_product models.Commodity_product, userId int) (modelName, error) {

	// 1. 取出对应产品
	var commodity modelName
	var product models.Product

	productRow := utils.DbQueryRow(db, "", "product", commodity_product.Product_id.Int, product)
	errGetProduct := product.ScanRow(productRow)

	if errGetProduct != nil {
		return commodity, errGetProduct
	}

	// 2. 用产品的属性来填充商品属性（作为初始值），并创建对应商品
	commodity.Name = product.Name
	commodity.Code = product.Code
	commodity.Memo = product.Memo
	commodity.Category_id = product.Category_id

	result, errInsert := utils.DbQueryInsert(db, tableName, commodity)

	if errInsert != nil {
		return commodity, errInsert
	}

	id, errId := result.LastInsertId()
	commodity.ID = nulls.NewInt(int(id))

	if errId != nil {
		return commodity, errId
	}

	// 3. 商品和产品进行绑定
	commodity_product.Commodity_id = commodity.ID
	commodity_product.IsMeta = nulls.NewInt(1) // 创建的时候是1，assemble的时候是0

	result, errInsert = utils.DbQueryInsert(db, "commodity_product", commodity_product)

	return commodity, errId
}

func (b repositoryName) GetRows_ByProduct(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	// 拦截 search.
	product_id := searchTerms["product_id"]
	delete(searchTerms, "product_id")

	// 根据product_id从中间表取非meta的
	subsql := fmt.Sprintf("(SELECT m1.*, m2.product_id as product_id, m3.image_id as image_id FROM "+tableName+" m1 LEFT JOIN commodity_product m2 ON m1.id = m2.commodity_id LEFT JOIN product m3 ON m2.product_id = m3.id WHERE m2.product_id = %s AND m2.isMeta = 0)", product_id)

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

func (b repositoryName) Assemble(db *sql.DB, commodity_id int, product_id int, userId int) error {

	_, err := db.Exec("INSERT INTO commodity_product (commodity_id, product_id, isMeta) VALUES(?, ?, 0);", commodity_id, product_id)

	var item modelName
	item.UpdateUser_id = nulls.NewInt(userId)
	item.ID = nulls.NewInt(commodity_id)
	_, err = utils.DbQueryUpdate(db, tableName, item)

	return err

}

func (b repositoryName) Disassemble(db *sql.DB, commodity_id int, product_id int, userId int) error {

	_, err := db.Exec("DELETE FROM commodity_product WHERE commodity_id = ? AND product_id=?;", commodity_id, product_id)

	var item modelName
	item.UpdateUser_id = nulls.NewInt(userId)
	item.ID = nulls.NewInt(commodity_id)
	_, err = utils.DbQueryUpdate(db, tableName, item)

	return err
}
