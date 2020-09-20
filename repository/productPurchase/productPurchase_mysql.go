package productPurchaseRepository

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.ProductPurchase
type repositoryName = Repository

var tableName = "product_purchase"

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

// =============================================== customized

func (b repositoryName) GetRows_GroupByCompany(
	mydb models.MyDb,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 product_id ；因为mysql的bug，去重取一的排序只能写在内部，没办法自动用searchTerm来做
	product_id := searchTerms["product_id"]
	// delete(searchTerms, "product_id")

	// rows, err := utils.DbQueryRows_Customized(mydb, "", tableName, &pagination, searchTerms, item, " GROUP BY company_id, currency_id")
	//select * from product_purchase where id in (select max(id) from product_purchase WHERE product_id = 10 group by company_id, currency_id) order by id desc
	// old: 	"(SELECT m1.* FROM product_purchase m1 LEFT JOIN product_purchase m2 ON (m1.company_id = m2.company_id AND m1.currency_id = m2.currency_id AND m1.id < m2.id) WHERE m2.id IS NULL)",

	sqlString := fmt.Sprintf(
		`(select * from product_purchase 
			where id in 
			(select max(id) from product_purchase 	
				WHERE product_id = %s 
				group by company_id, currency_id, memo) 
				order by id desc)`, product_id)

	rows, err := utils.DbQueryRows_Customized(mydb, "",
		sqlString,
		&pagination,
		searchTerms,
		item,
		"",
		"")

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {

		err = item.ScanRows(rows)

		if err != nil {
			return []modelName{}, pagination, err
		}

		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRows_History(
	mydb models.MyDb,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 search
	productpurchase_id := searchTerms["productpurchase_id"]
	delete(searchTerms, "productpurchase_id")

	// var sqlString string

	// root_id_int, err := strconv.Atoi(root_id)

	// rows, err := utils.DbQueryRows_Customized(mydb, "", tableName, &pagination, searchTerms, item, " GROUP BY company_id, currency_id")
	rows, err := utils.DbQueryRows_Customized(mydb, "",
		"(SELECT a.* FROM product_purchase a LEFT JOIN product_purchase b ON a.company_id = b.company_id AND a.product_id = b.product_id WHERE b.id = "+productpurchase_id+")",
		&pagination,
		searchTerms,
		item,
		"",
		"")

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {

		err = item.ScanRows(rows)

		if err != nil {
			return []modelName{}, pagination, err
		}

		items = append(items, item)
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow_ByProductId(mydb models.MyDb, product_id int, company_id int, userId int) (modelName, error) {

	var item modelName

	var row *sql.Row

	if company_id > 0 {
		query := "SELECT a.* FROM " + tableName + " a WHERE a.product_id = ? AND company_id = ? ORDER BY a.activeAt DESC LIMIT 1"

		if mydb.Tx != nil {
			row = mydb.Tx.QueryRow(query, product_id, company_id)
		} else {
			row = mydb.Db.QueryRow(query, product_id, company_id)
		}

	} else {
		query := "SELECT a.* FROM " + tableName + " a WHERE a.product_id = ?  ORDER BY a.activeAt DESC LIMIT 1"
		if mydb.Tx != nil {
			row = mydb.Tx.QueryRow(query, product_id)
		} else {
			row = mydb.Db.QueryRow(query, product_id)
		}

	}

	err := row.Scan(item.Receivers()...)

	return item, err
}
