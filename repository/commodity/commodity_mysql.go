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
// 还是判断meta。删了产品对应的商品和对应的关联也应该删掉了:
// 		删除产品的时候，判断commodity_product有没有isMeta的。如果有就删除商品. 然后删除所有关联产品的 commodity_product
// 		删除商品的时候也删除所有对应的commodity_product

var tableMeta = `(
		SELECT 
			core1.*, 
			core2.product_id as product_id, 
			core3.image_id as image_id 
			FROM commodity core1 
				LEFT JOIN commodity_product core2  
					ON core1.id = core2.commodity_id AND core2.isMeta = 1
				LEFT JOIN product core3 
					ON core3.id = core2.product_id)`

var tableAll = `(
	SELECT 
		core1.*, 
		core2.product_id as product_id, 
		core3.image_id as image_id 
		FROM commodity core1 
			LEFT JOIN commodity_product core2  
				ON core1.id = core2.commodity_id
			LEFT JOIN product core3 
				ON core3.id = core2.product_id)`

/*
	产品和meta应该是一一绑定的。只要有商品，就一定有一个meta。而组合商品作为极少数情况，额外挂载
		所以是两种情况：meta就是“产品+商品”的单条，否则就是挂载
		删除产品，meta的商品也删除，对应的挂载就全部删除（组合商品）: 全删除的目前用cascade
*/

// =============================================== basic CRUD
func (b repositoryName) GetRows(
	db *sql.DB,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

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

func (b repositoryName) GetRow(db *sql.DB, id int, userId int) (modelName, error) {
	var item modelName

	row := utils.DbQueryRow(db, "", tableMeta, id, item)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) GetRow_ByProduct(db *sql.DB, id int, userId int) (modelName, error) {
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

	// 数据库没有这个，这是为了显示用的，所以去掉
	item.Image_id = nulls.Int{Int: 0, Valid: false}

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

	// 数据库没有这个，这是为了显示用的，所以去掉
	item.Image_id = nulls.Int{Int: 0, Valid: false}

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

	// result, row, err := utils.DbQueryDelete(db, tableName, id, item)
	// 这里特殊，返回一个空的item。因为commodity的image_id字段是在关联的product表里，这里delete的地方取不到，而且外部用不到这个item
	result, row, err := utils.DbQueryDelete(db, tableName, tableName, id, item)
	item.ScanRow(row)

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

func (b repositoryName) GetPrintSource(db *sql.DB, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(db, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
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
	commodity.EName = product.EName
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
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 search.
	product_id := searchTerms["product_id"]
	delete(searchTerms, "product_id")

	// 根据product_id从中间表取(为避免用户迷惑，meta也一起取)
	subsql := fmt.Sprintf("(SELECT m1.*, m2.product_id as product_id, m3.image_id as image_id FROM "+tableName+" m1 LEFT JOIN commodity_product m2 ON m1.id = m2.commodity_id LEFT JOIN product m3 ON m2.product_id = m3.id WHERE m2.product_id = %s)", product_id)

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

	result, err := db.Exec("INSERT INTO commodity_product (commodity_id, product_id, isMeta) VALUES(?, ?, 0);", commodity_id, product_id)
	rowsUpdated, err := result.RowsAffected()

	if rowsUpdated > 0 {
		var item modelName
		item.UpdateUser_id = nulls.NewInt(userId)
		item.ID = nulls.NewInt(commodity_id)
		_, row, errUpdate := utils.DbQueryUpdate(db, tableName, tableName, item)
		item.ScanRow(row)
		return errUpdate
	}

	return err
}

func (b repositoryName) Disassemble(db *sql.DB, commodity_id int, product_id int, userId int) error {

	// 看有没有包含主产品。如果有的话删掉商品信息
	_, err := db.Exec(`
	DELETE a FROM commodity a JOIN commodity_product b
	ON a.id = b.commodity_id	
	WHERE b.commodity_id = ? AND b.product_id=? AND isMeta = 1;
	`, commodity_id, product_id)

	// 删掉关联
	result, err := db.Exec("DELETE FROM commodity_product WHERE commodity_id = ? AND product_id=?;", commodity_id, product_id)
	rowsUpdated, err := result.RowsAffected()

	if rowsUpdated > 0 {
		var item modelName
		item.UpdateUser_id = nulls.NewInt(userId)
		item.ID = nulls.NewInt(commodity_id)
		_, row, errUpdate := utils.DbQueryUpdate(db, tableName, tableName, item)
		item.ScanRow(row)
		return errUpdate
	}

	return err
}
