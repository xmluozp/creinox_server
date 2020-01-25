package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Folder struct {
	ID         int          `col:"" json:"id"`
	Memo       nulls.String `col:"" json:"memo"`
	FolderType nulls.Int    `col:"" json:"FolderType"` // 权限可以后期改
	ViewSource nulls.String `col:"" json:"ViewSource"`
	RefSource  nulls.String `col:"" json:"refSource"`
	RefId      nulls.Int    `col:"" json:"RefId"`
}

type FolderList struct {
	Items []*Folder
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *Folder) ScanRow(r *sql.Row) error {
	return r.Scan(
		&item.ID,
		&item.Memo,
		&item.FolderType,
		&item.ViewSource,
		&item.RefSource,
		&item.RefId)
}

func (item *Folder) ScanRows(r *sql.Rows) error {
	return r.Scan(
		&item.ID,
		&item.Memo,
		&item.FolderType,
		&item.ViewSource,
		&item.RefSource,
		&item.RefId)
}

func (list *FolderList) ScanRow(r *sql.Rows) error {

	item := new(Folder) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
