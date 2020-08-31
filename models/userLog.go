package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

// TODO: Snapshot怎么办？ 新的snapshot可以从前台取
type UserLog struct {
	ID             nulls.Int    `col:"" json:"id"`
	CreateAt       nulls.Time   `col:"" json:"createAt"`
	Memo           nulls.String `col:"" json:"memo"`
	Type           nulls.Int    `col:"" json:"type"`
	SnapshotBefore nulls.String `col:"" json:"snapshotBefore"`
	SnapshotAfter  nulls.String `col:"" json:"snapshotAfter"`
	UpdateUser_id  nulls.Int    `col:"fk" json:"updateUser_id"`

	//========fk
	UpdateUser User `ref:"user,updateUser_id" json:"updateUser_id.row" validate:"-"`
}

func (item *UserLog) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.CreateAt,
		&item.Memo,
		&item.Type,
		&item.SnapshotBefore,
		&item.SnapshotAfter,
		&item.UpdateUser_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *UserLog) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkUpdateUser := User{}

	columns = append(item.Receivers(), fkUpdateUser.Receivers()...)

	err := r.Scan(columns...)

	item.UpdateUser = fkUpdateUser

	return err
}

func (item *UserLog) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkUpdateUser := User{}

	columns = append(item.Receivers(), fkUpdateUser.Receivers()...)

	err := r.Scan(columns...)

	item.UpdateUser = fkUpdateUser

	return err
}
