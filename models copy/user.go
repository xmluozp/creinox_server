package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

// validator参考：https://github.com/go-playground/validator/blob/v9/_examples/simple/main.go
type User struct {
	ID          nulls.Int    `col:"" json:"id"`
	UserName    nulls.String `col:"" json:"userName" validate:"required" errm:"用户名必填"`
	FullName    nulls.String `col:"" json:"fullName" `
	Password    nulls.String `col:"" json:"password" validate:"omitempty,gt=3,lt=16" errm:"密码是3-16位长的字符"`
	IP          nulls.String `col:"" json:"ip"`
	BargainCode nulls.String ` json:"bargainCode,omitempty"`
	LastLogin   nulls.Time   `col:"" json:"lastLogin"` // 先试一下，不行就改用string. 读取：time.Parse("2006/01/02", ranges[0])
	CreateAt    nulls.Time   `col:"default" json:"createAt"`
	Token       nulls.String `col:"" json:"token"`
	Memo        nulls.String `col:"" json:"memo"`
	IsActive    nulls.Bool   `col:"" json:"isActive"`
	Role_id     nulls.Int    `col:"fk" json:"role_id" validate:"required" errm:"角色必选"`

	RoleItem Role `ref:"role,role_id" json:"role_id.row" validate:"-"`
}

type UserList struct {
	Items []*User
}

func (item *User) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.UserName,
		&item.FullName,
		&item.Password,
		&item.IP,
		&item.BargainCode,
		&item.LastLogin,
		&item.CreateAt,
		&item.Token,
		&item.Memo,
		&item.IsActive,
		&item.Role_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *User) ScanRow(r *sql.Row) error {

	var columns []interface{}

	RoleItem := Role{}

	columns = append(item.Receivers(), RoleItem.Receivers()...)

	err := r.Scan(columns...)

	item.RoleItem = RoleItem

	return err
}

func (item *User) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	RoleItem := Role{}

	columns = append(item.Receivers(), RoleItem.Receivers()...)

	err := r.Scan(columns...)

	item.RoleItem = RoleItem

	return err
}

func (list *UserList) ScanRow(r *sql.Rows) error {

	item := new(User) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
