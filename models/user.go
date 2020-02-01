package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

// validator参考：https://github.com/go-playground/validator/blob/v9/_examples/simple/main.go
type User struct {
	ID          int          `col:"" json:"id"`
	UserName    nulls.String `col:"" json:"userName" validate:"required" errm:"用户名必填"`
	FullName    nulls.String `col:"" json:"fullName" `
	Password    nulls.String `col:"" json:"password" validate:"gt=3,lt=16" errm:"密码是3-16位长的字符"`
	IP          nulls.String `col:"" json:"ip"`
	BargainCode nulls.String ` json:"bargainCode,omitempty"`
	LastLogin   nulls.Time   `col:"" json:"lastLogin"` // 先试一下，不行就改用string. 读取：time.Parse("2006/01/02", ranges[0])
	CreateAt    nulls.Time   `col:"default" json:"createAt"`
	Token       nulls.String `col:"" json:"token"`
	Memo        nulls.String `col:"" json:"memo"`
	IsActive    nulls.Bool   `col:"" json:"isActive"`
	Role_id     nulls.Int    `col:"fk" json:"role_id" validate:"required" errm:"角色必选"`
}

type UserList struct {
	Items []*User
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *User) ScanRow(r *sql.Row) error {
	return r.Scan(
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
		&item.Role_id)
}

func (item *User) ScanRows(r *sql.Rows) error {
	return r.Scan(
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
		&item.Role_id)
}

func (list *UserList) ScanRow(r *sql.Rows) error {

	item := new(User) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
