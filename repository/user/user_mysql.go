package userRepository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Repository struct{}
type modelName = models.User
type repositoryName = Repository

var tableName = "user"

// =============================================== Login
func (b repositoryName) GetLoginRow(db *sql.DB, userName string) (modelName, error) {

	var item modelName

	// passed in is the encryped password
	row := db.QueryRow("SELECT * FROM "+tableName+" WHERE userName = ? AND isActive = TRUE", userName)

	err := item.ScanRow(row)

	return item, err
}

func (b repositoryName) UpdateLoginRow(db *sql.DB, item modelName) (int64, error) {

	// login更新，只更新ip, 上次登录，token
	result, err := db.Exec("UPDATE user SET ip=?, lastLogin = CURRENT_TIMESTAMP, token=? WHERE id=?", &item.IP, &item.Token, &item.ID)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

// =============================================== basic CRUD

func (b repositoryName) GetRows(
	db *sql.DB,
	item modelName,
	items []modelName,
	pagination models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]modelName, models.Pagination, error) {

	// rows这里是一个cursor.
	rows, err := utils.DbQueryRows(db, "SELECT * FROM "+tableName+" WHERE 1=1 ", tableName, &pagination, searchTerms, item)

	if err != nil {
		return []modelName{}, pagination, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {
		err = item.ScanRows(rows)
		items = append(items, item)

		if err != nil {
			fmt.Println("rows:", rows, err)
		}
	}

	if err != nil {
		return []modelName{}, pagination, err
	}

	return items, pagination, nil
}

func (b repositoryName) GetRow(db *sql.DB, id int) (modelName, error) {
	var item modelName
	row := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	err := item.ScanRow(row)

	// 不显示密码
	item.Password = nulls.NewString("")

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	hashedPass, _ := auth.HashPassword(item.Password.String)
	item.Password = nulls.NewString(hashedPass)

	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = int(id)
	if errId != nil {
		return item, errId
	}

	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	var result sql.Result
	var err error

	if item.Password.String != "" {
		hashedPass, _ := auth.HashPassword(item.Password.String)
		item.Password = nulls.NewString(hashedPass)
		result, err = utils.DbQueryUpdate(db, tableName, item)
	} else {
		result, err = db.Exec("UPDATE user SET fullName = ?, memo = ?, isActive = ?, role_id=? WHERE id=?", &item.FullName, &item.Memo, &item.IsActive, &item.Role_id, &item.ID)
	}

	// result, err := db.Exec("UPDATE role SET name=?, rank = ?, auth=? WHERE id=?", &item.Name, &item.Rank, &item.Auth, &item.ID)

	result, err = utils.DbQueryUpdate(db, tableName, item)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

func (b repositoryName) DeleteRow(db *sql.DB, id int, userId int) (int64, error) {

	result, err := db.Exec("DELETE FROM "+tableName+" WHERE id = ?", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, err
}
