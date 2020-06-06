package userRepository

import (
	"database/sql"
	"errors"
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
	row := db.QueryRow("SELECT * FROM "+tableName+" WHERE userName = ? AND isActive = 1", userName)

	err := row.Scan(item.Receivers()...)
	// item.ScanRow(row)

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
	searchTerms map[string]string,
	userId int) ([]modelName, models.Pagination, error) {

	rank := auth.GetRankFromUser(db, userId)

	// 搜索roles比自己小的
	subsql := fmt.Sprintf("(SELECT a.* FROM user a LEFT JOIN role b ON b.id = a.role_id WHERE b.rank > %d OR b.rank IS NULL OR a.id = %d)", rank, userId)

	rows, err := utils.DbQueryRows(db, "", subsql, &pagination, searchTerms, item)

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

func (b repositoryName) GetRow(db *sql.DB, id int, userId int) (modelName, error) {
	var item modelName
	row := utils.DbQueryRow(db, "", tableName, id, item)

	err := item.ScanRow(row)

	// 不显示密码
	item.Password = nulls.NewString("")

	return item, err
}

func (b repositoryName) AddRow(db *sql.DB, item modelName, userId int) (modelName, error) {

	// 判断用户名是否唯一

	count := 0
	scanErr := db.QueryRow("SELECT COUNT(*) FROM " + tableName + " WHERE userName = '" + item.UserName.String + "'").Scan(&count)

	if scanErr != nil {
		return item, scanErr
	}

	if count > 0 {
		return item, errors.New(" 用户名已存在")
	}

	//
	hashedPass, _ := auth.HashPassword(item.Password.String)
	item.Password = nulls.NewString(hashedPass)

	result, errInsert := utils.DbQueryInsert(db, tableName, item)

	if errInsert != nil {
		return item, errInsert
	}

	id, errId := result.LastInsertId()
	item.ID = nulls.NewInt(int(id))
	if errId != nil {
		return item, errId
	}

	item.Password = nulls.NewString("")
	return item, errId
}

func (b repositoryName) UpdateRow(db *sql.DB, item modelName, userId int) (int64, error) {

	// 判断用户名是否唯一

	count := 0
	scanErr := db.QueryRow("SELECT COUNT(*) FROM " + tableName + " WHERE userName = '" + item.UserName.String + "'").Scan(&count)

	if scanErr != nil {
		return 0, scanErr
	}

	if count > 0 {
		return 0, errors.New(" 用户名已存在")
	}

	//
	var result sql.Result
	var err error

	if item.Password.String != "" {
		hashedPass, _ := auth.HashPassword(item.Password.String)
		item.Password = nulls.NewString(hashedPass)
		result, _, err = utils.DbQueryUpdate(db, tableName, tableName, item)
	} else {
		// 防止最高管理员把自己禁用或者降级
		if item.ID.Int == userId {
			result, err = db.Exec("UPDATE user SET fullName = ?, memo = ? WHERE id=?", &item.FullName, &item.Memo, &item.ID)

		} else {
			result, err = db.Exec("UPDATE user SET fullName = ?, memo = ?, isActive = ?, role_id=? WHERE id=?", &item.FullName, &item.Memo, &item.IsActive, &item.Role_id, &item.ID)
		}
	}

	// result, err := db.Exec("UPDATE role SET name=?, rank = ?, auth=? WHERE id=?", &item.Name, &item.Rank, &item.Auth, &item.ID)

	// result, err = utils.DbQueryUpdate(db, tableName, item)

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

	if id == userId {
		return nil, errors.New("You can not delete yourself")
	}

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

	return item, err
}

func (b repositoryName) GetPrintSource(db *sql.DB, id int, userId int) (modelName, error) {

	return b.GetRow(db, id, userId)
}
