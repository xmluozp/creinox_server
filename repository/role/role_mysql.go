package roleRepository

import (
	"database/sql"
	"github.com/xmluozp/creinox_server/models"
)

type RoleRepository struct{}

func (b RoleRepository) GetRoles(db *sql.DB, role models.Role, roles []models.Role) ([]models.Role, error) {
	// rows这里是一个cursor
	rows, err := db.Query("select * from roles")

	if err != nil {
		return []models.Role{}, err
	}

	defer rows.Close() // 以下代码执行完了，关闭连接

	for rows.Next() {

		// 把数据库读出来的列填进对应的变量里 (如果只想取对应的列怎么办？)
		err = rows.Scan(&role.ID, &role.Name, &role.Rank, &role.Auth)
		roles = append(roles, role)
	}

	if err != nil {
		return []models.Role{}, err
	}

	return roles, nil
}

// ?? 为什么要用外来的变量，却传回一个拷贝？？
func (b RoleRepository) GetRole(db *sql.DB, id int) (models.Role, error) {
	var role models.Role
	rows := db.QueryRow("SELECT * FROM roles WHERE id = $1", id)

	// 假如不是平的struct而有子选项
	// 就要改写Scan
	// https://stackoverflow.com/questions/47335697/golang-decode-json-request-in-nested-struct-and-insert-in-db-as-blob
	err := rows.Scan(&role.ID, &role.Name, &role.Rank, &role.Auth)

	return role, err
}

func (b RoleRepository) AddRole(db *sql.DB, role models.Role) (int, error) {

	err := db.QueryRow("insert into roles (title, author, year) values($1, $2, $3) returning id;",
		role.Name, role.Rank, role.Auth).Scan(&role.ID)

	if err != nil {
		return 0, err
	}

	return role.ID, err
}

func (b RoleRepository) UpdateRole(db *sql.DB, role models.Role) (int64, error) {

	result, err := db.Exec("update roles set title=$1, author = $2, year=$3 where id=$4 RETURNING id",
		&role.Name, &role.Rank, &role.Auth, &role.ID)

	if err != nil {
		return 0, err
	}
	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}
	return rowsUpdated, err
}

func (b RoleRepository) DeleteRole(db *sql.DB, id int) (int64, error) {

	result, err := db.Exec("delete from roles where id = $1", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, err
}
