package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type CommonItem struct {
	ID         nulls.Int    `col:"" json:"id"`
	Name       nulls.String `col:"" json:"name" validate:"required" errm:"必填"`
	Ename      nulls.String `col:"" json:"ename"`
	Memo       nulls.String `col:"" json:"memo"`
	Auth       nulls.String `col:"" json:"auth"`
	Sorting    nulls.Int    `col:"" json:"sorting"`
	IsActive   nulls.Bool   `col:"" json:"isActive"`
	IsDelete   nulls.Bool   `col:"default" json:"isDelete"`
	CommonType nulls.Int    `col:"" json:"commonType"`
}

type CommonItemList struct {
	Items []*CommonItem
}

func (item *CommonItem) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Name,
		&item.Ename,
		&item.Memo,
		&item.Auth,
		&item.Sorting,
		&item.IsActive,
		&item.IsDelete,
		&item.CommonType}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *CommonItem) ScanRow(r *sql.Row) error {
	return r.Scan(item.Receivers()...)
}

func (item *CommonItem) ScanRows(r *sql.Rows) error {
	return r.Scan(item.Receivers()...)
}

func (list *CommonItemList) ScanRow(r *sql.Rows) error {

	item := new(CommonItem) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
