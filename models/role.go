package models

import (
	"database/sql"
	"github.com/gobuffalo/nulls"
)

type Role struct {
	ID   int          `col:"" json:"id"`
	Name nulls.String `col:"" json:"name" validate:"required" errm:"角色名必填"`
	Rank nulls.Int    `col:"" json:"rank"` // 权限可以后期改
	Auth nulls.String `col:"" json:"auth" validate:"required" errm:"必须选择"`
}

type RoleList struct {
	Items []*Role
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (role *Role) ScanRow(r *sql.Row) error {
	return r.Scan(&role.ID, &role.Name, &role.Rank, &role.Auth)
}

func (role *Role) ScanRows(r *sql.Rows) error {
	return r.Scan(&role.ID, &role.Name, &role.Rank, &role.Auth)
}

func (list *RoleList) ScanRow(r *sql.Rows) error {

	item := new(Role) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
