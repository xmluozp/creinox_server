package financialLedgerRepository

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.FinancialLedger
type repositoryName = Repository

var tableName = "financial_ledger"

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

	// SELECT * FROM financialLedger WHERE path LIKE CONCAT((SELECT path FROM financialLedger WHERE id = 1), ',',1, '%') ORDER BY path ASC
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
	row := utils.DbQueryRow(mydb, "", tableName, id, item)

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

	err := b.updateLedgerNames(mydb, item.ID.Int)

	return item, err
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

	err = b.updateLedgerNames(mydb, item.ID.Int)

	return rowsUpdated, err
}
func (b repositoryName) updateLedgerNames(mydb models.MyDb, id int) (err error) {

	// 更新下属节点的name
	var updateQuery = `
		UPDATE financial_ledger main LEFT JOIN    
		(
			SELECT a.id, (
				SELECT CONCAT(GROUP_CONCAT(name SEPARATOR '/'), '/', a.name) FROM financial_ledger b
				WHERE find_in_set(b.id, a.path)
			) AS newLedgerName FROM financial_ledger a) sub
		ON main.id = sub.id JOIN
		(
			SELECT path, id FROM financial_ledger WHERE id = %d
		) conditions
		SET main.ledgerName = sub.newLedgerName
		WHERE 
		main.id = conditions.id OR
		main.path = CONCAT(conditions.path, ',' , conditions.id) OR
		main.path LIKE CONCAT(conditions.path, ',' , conditions.id, ',', '%%')
	`
	updateQueryCombined := fmt.Sprintf(updateQuery, id)

	if mydb.Tx != nil {
		_, err = mydb.Tx.Exec(updateQueryCombined)
	} else {
		_, err = mydb.Db.Exec(updateQueryCombined)
	}

	return err
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
