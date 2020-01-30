package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Region struct {
	ID        int          `col:"" json:"id"`
	Name      nulls.String `col:"" json:"name" validate:"required" errm:"地区名必填"`
	Ename     nulls.String `col:"" json:"ename"`
	TelPrefix nulls.String `col:"" json:"telPrefix"`
	Code      nulls.String `col:"" json:"code"`
	TreeLock  nulls.Bool   `col:"" json:"treeLock"`
	Memo      nulls.String `col:"" json:"memo"`
	Path      nulls.String `col:"" json:"path"`
	UpdateAt  nulls.Time   `col:"newtime" json:"updateAt"`
	CreateAt  nulls.Time   `col:"default" json:"createAt"`
	IsDelete  nulls.Bool   `col:"" json:"isDelete"`
	Parent_id nulls.Int    `col:"" json:"parent_id"`
	Root_id   nulls.Int    `json:"root_id"`
}

type RegionList struct {
	Items []*Region
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *Region) ScanRow(r *sql.Row) error {
	return r.Scan(
		&item.ID,
		&item.Name,
		&item.Ename,
		&item.TelPrefix,
		&item.Code,
		&item.TreeLock,
		&item.Memo,
		&item.Path,
		&item.UpdateAt,
		&item.CreateAt,
		&item.IsDelete,
		&item.Parent_id)
}

func (item *Region) ScanRows(r *sql.Rows) error {
	return r.Scan(
		&item.ID,
		&item.Name,
		&item.Ename,
		&item.TelPrefix,
		&item.Code,
		&item.TreeLock,
		&item.Memo,
		&item.Path,
		&item.UpdateAt,
		&item.CreateAt,
		&item.IsDelete,
		&item.Parent_id)
}

func (list *RegionList) ScanRow(r *sql.Rows) error {

	item := new(Region) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
