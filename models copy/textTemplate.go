package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type TextTemplate struct {
	ID          nulls.Int    `col:"" json:"id"`
	Name        nulls.String `col:"" json:"name"`
	TargetTable nulls.String `col:"" json:"targetTable"  validate:"required" errm:"目标表格必填"`
	ColumnName  nulls.String `col:"" json:"columnName"  validate:"required" errm:"目标列必填"`
	Content     nulls.String `col:"" json:"content"`
	UpdateAt    nulls.Time   `col:"newtime" json:"updateAt"`
}

type TextTemplateList struct {
	Items []*TextTemplate
}

func (item *TextTemplate) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID, &item.Name, &item.TargetTable, &item.ColumnName, &item.Content, &item.UpdateAt}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *TextTemplate) ScanRow(r *sql.Row) error {

	err := r.Scan(item.Receivers()...)

	return err
}

func (item *TextTemplate) ScanRows(r *sql.Rows) error {
	err := r.Scan(item.Receivers()...)

	return err
}

func (list *TextTemplateList) ScanRow(r *sql.Rows) error {

	item := new(TextTemplate) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
