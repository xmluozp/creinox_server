package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Category struct {
	ID          nulls.Int    `col:"" json:"id"`
	Name        nulls.String `col:"" json:"name" validate:"required" errm:"必填"`
	Ename       nulls.String `col:"" json:"ename"`
	Prefix      nulls.String `col:"" json:"prefix"`
	CurrentCode nulls.String `col:"" json:"currentCode"`
	TreeLock    nulls.Bool   `col:"" json:"treeLock"`
	Memo        nulls.String `col:"" json:"memo"`
	Ememo       nulls.String `col:"" json:"ememo"`
	Path        nulls.String `col:"" json:"path"`
	UpdateAt    nulls.Time   `col:"newtime" json:"updateAt"`
	CreateAt    nulls.Time   `col:"default" json:"createAt"`
	IsDelete    nulls.Bool   `col:"default" json:"isDelete"`
	Parent_id   nulls.Int    `col:"fk" json:"parent_id"`
	Root_id     nulls.Int    `json:"root_id"` // to display tree
}

type CategoryList struct {
	Items []*Category
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *Category) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Name,
		&item.Ename,
		&item.Prefix,
		&item.CurrentCode,
		&item.TreeLock,
		&item.Memo,
		&item.Ememo,
		&item.Path,
		&item.UpdateAt,
		&item.CreateAt,
		&item.IsDelete,
		&item.Parent_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *Category) ScanRow(r *sql.Row) error {
	return r.Scan(item.Receivers()...)
}

func (item *Category) ScanRows(r *sql.Rows) error {
	return r.Scan(item.Receivers()...)
}

func (list *CategoryList) ScanRow(r *sql.Rows) error {

	item := new(Category) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
