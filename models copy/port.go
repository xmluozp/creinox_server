package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Port struct {
	ID            nulls.Int    `col:"" json:"id"`
	Name          nulls.String `col:"" json:"name" validate:"required" errm:"名字必填"`
	EName         nulls.String `col:"" json:"ename"`
	IsDeparture   nulls.Bool   `col:"" json:"isDeparture"` // 权限可以后期改
	IsDestination nulls.Bool   `col:"" json:"isDestination"`
}

type PortList struct {
	Items []*Port
}

func (item *Port) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID, &item.Name, &item.EName, &item.IsDeparture, &item.IsDestination}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *Port) ScanRow(r *sql.Row) error {

	err := r.Scan(item.Receivers()...)

	return err
}

func (item *Port) ScanRows(r *sql.Rows) error {
	err := r.Scan(item.Receivers()...)

	return err
}

func (list *PortList) ScanRow(r *sql.Rows) error {

	item := new(Port) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
