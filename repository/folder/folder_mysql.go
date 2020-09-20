package folderRepository

import (
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.Folder
type repositoryName = Repository

var tableName = "folder"

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

func (b repositoryName) AddRow_withRef(mydb models.MyDb, item modelName, userId int) (returnItem modelName, err error) {

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
	// 更新源

	sqlstr := fmt.Sprintf("UPDATE %s SET %s=? WHERE id=?", item.TableName.String, item.ColumnName.String)

	if mydb.Tx != nil {
		_, err = mydb.Tx.Exec(sqlstr, item.ID.Int, item.RefId.Int)
	} else {
		_, err = mydb.Db.Exec(sqlstr, item.ID.Int, item.RefId.Int)
	}

	if err != nil {
		return item, err
	}

	return item, errId
}

func (b repositoryName) GetPrintSource(mydb models.MyDb, id int, userId int) (map[string]interface{}, error) {

	item, err := b.GetRow(mydb, id, userId)

	if err != nil {
		return nil, err
	}

	ds, err := utils.GetPrintSourceFromInterface(item)

	return ds, err
}
