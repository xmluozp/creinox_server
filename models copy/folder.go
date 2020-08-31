package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Folder struct {
	ID         nulls.Int    `col:"" json:"id"`
	Memo       nulls.String `col:"" json:"memo"`
	FolderType nulls.Int    `col:"" json:"folderType"` // 权限可以后期改
	ViewSource nulls.String `col:"" json:"viewSource"`
	RefSource  nulls.String `col:"" json:"refSource"`
	RefId      nulls.Int    `col:"" json:"refId"`

	// 生成folder的时候用来返插入源表格的
	TableName  nulls.String `json:"tableName"`
	ColumnName nulls.String `json:"columnName"`
}

type FolderList struct {
	Items []*Folder
}

func (item *Folder) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Memo,
		&item.FolderType,
		&item.ViewSource,
		&item.RefSource,
		&item.RefId}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *Folder) ScanRow(r *sql.Row) error {

	err := r.Scan(item.Receivers()...)

	return err
}

func (item *Folder) ScanRows(r *sql.Rows) error {

	err := r.Scan(item.Receivers()...)

	return err
}

func (list *FolderList) ScanRow(r *sql.Rows) error {

	item := new(Folder) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
