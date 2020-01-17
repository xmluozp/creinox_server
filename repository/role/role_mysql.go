package roleRepository

import (
	"database/sql"
	"fmt"

	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type RoleRepository struct{}

func (b RoleRepository) GetRoles(
	db *sql.DB,
	role models.Role,
	roles []models.Role,
	pagination *models.Pagination, // 需要返回总页数
	searchTerms map[string]string) ([]models.Role, error) {

	// rows这里是一个cursor
	rows, err := utils.DbQueryRows(db, "select * from role", "role", pagination, searchTerms, role)

	if err != nil {
		return []models.Role{}, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {

		// 把数据库读出来的列填进对应的变量里 (如果只想取对应的列怎么办？)
		// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

		rows.Scan(&role.ID, &role.Name, &role.Rank, &role.Auth)
		roles = append(roles, role)
	}

	if err != nil {
		return []models.Role{}, err
	}

	// for i, _ := range roles {
	// 	roles[i].Name = "改头换面"
	// }

	return roles, nil
}

// ?? 为什么要用外来的变量，却传回一个拷贝？？
func (b RoleRepository) GetRole(db *sql.DB, id int) (models.Role, error) {
	var role models.Role
	rows := db.QueryRow("SELECT * FROM role WHERE id = ?", id)

	// 假如不是平的struct而有子选项
	// 就要改写Scan
	// https://stackoverflow.com/questions/47335697/golang-decode-json-request-in-nested-struct-and-insert-in-db-as-blob
	err := rows.Scan(&role.ID, &role.Name, &role.Rank, &role.Auth)

	return role, err
}

func (b RoleRepository) AddRole(db *sql.DB, role models.Role) (models.Role, error) {

	result, errInsert := db.Exec("insert into role (name, rank, auth) values(?, ?, ?);", role.Name, role.Rank, role.Auth)
	if errInsert != nil {
		fmt.Println("insert error: ", errInsert)
		return role, errInsert
	}

	id, errId := result.LastInsertId()
	role.ID = int(id)
	if errId != nil {
		return role, errId
	}

	return role, errId
}

func (b RoleRepository) UpdateRole(db *sql.DB, role models.Role) (int64, error) {

	result, err := db.Exec("UPDATE role SET name=?, rank = ?, auth=? WHERE id=?", &role.Name, &role.Rank, &role.Auth, &role.ID)

	if err != nil {
		return 0, err
		fmt.Println("update error: ", err)
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, err
}

func (b RoleRepository) DeleteRole(db *sql.DB, id int) (int64, error) {

	result, err := db.Exec("DELETE FROM role WHERE id = ?", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, err
}
