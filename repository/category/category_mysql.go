package categoryRepository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Category
type repositoryName = Repository

var tableName = "category"

// =============================================== basic CRUD

func (b repositoryName) GetRows(
	mydb models.MyDb,
	pagination models.Pagination,
	searchTerms map[string]string,
	userId int) (items []modelName, returnPagination models.Pagination, err error) {
	var item modelName

	// 拦截 search
	root_id := searchTerms["root_id"]
	delete(searchTerms, "root_id")

	var subsql string

	root_id_int, err := strconv.Atoi(root_id)

	if err == nil && root_id_int > 0 {
		subsql = fmt.Sprintf(
			`SELECT * FROM %s a JOIN (
			SELECT path, id FROM %s WHERE id =%d) b
			WHERE a.path = CONCAT(b.path, ',', b.id) or
			a.path LIKE CONCAT(b.path, ',', b.id, ',', '%%')`, tableName, tableName, root_id_int)
	} else {
		subsql = ""
	}
	rows, err := utils.DbQueryRows(mydb, subsql, tableName, &pagination, searchTerms, item)

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

	var row *sql.Row

	if mydb.Tx != nil {
		row = mydb.Tx.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)
	} else {
		row = mydb.Db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)
	}

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) AddRow(mydb models.MyDb, item modelName, userId int) (modelName, error) {

	// result, errInsert := db.Exec("INSERT INTO role (name, rank, auth) VALUES(?, ?, ?);", item.Name, item.Rank, item.Auth)
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

// ----------------------------------
func (b repositoryName) UpdateRow_currentCode(mydb models.MyDb, id int, code string, userId int) (int64, error) {

	var item modelName
	fmt.Println("获取分类当前号码")
	item, err := b.GetRow(mydb, id, userId)

	item.CurrentCode = nulls.NewString(code)
	item.ProductCount = nulls.NewInt(item.ProductCount.Int + 1)

	result, row, err := utils.DbQueryUpdate(mydb, tableName, tableName, item)

	// 之后改成重新统计
	// UPDATE `creinox`.`category` a
	// LEFT JOIN (SELECT COUNT(*) pc, category_id  FROM `creinox`.`product` GROUP BY category_id) b
	// ON a.id = b.category_id
	// SET a.productCount = b.pc;

	item.ScanRow(row)

	if err != nil {

		fmt.Println("获取分类当前号码出错")
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

func (b repositoryName) GetPrintSource(mydb models.MyDb, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
}
