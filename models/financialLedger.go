package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type FinancialLedger struct {
	ID        nulls.Int    `col:"" json:"id"`
	Name      nulls.String `col:"" json:"name" validate:"required" errm:"科目名必填"`
	Code      nulls.String `col:"" json:"code"`
	IsBuiltin nulls.Bool   `col:"" json:"isBuiltin"`
	Memo      nulls.String `col:"" json:"memo"`
	Path      nulls.String `col:"" json:"path"`
	Auth      nulls.String `col:"" json:"auth"`
	IsActive  nulls.Bool   `col:"" json:"isActive"`
	Sorting   nulls.Int    `col:"orderByAsc" json:"sorting"`
	Parent_id nulls.Int    `col:"fk" json:"parent_id"`
	Root_id   nulls.Int    `json:"root_id"`
}

type FinancialLedgerList struct {
	Items []*FinancialLedger
}

func (item *FinancialLedger) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Name,
		&item.Code,
		&item.IsBuiltin,
		&item.Memo,
		&item.Path,
		&item.Auth,
		&item.IsActive,
		&item.Sorting,
		&item.Parent_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *FinancialLedger) ScanRow(r *sql.Row) error {
	err := r.Scan(item.Receivers()...)
	return err
}

func (item *FinancialLedger) ScanRows(r *sql.Rows) error {
	err := r.Scan(item.Receivers()...)
	return err
}

func (list *FinancialLedgerList) ScanRow(r *sql.Rows) error {

	item := new(FinancialLedger) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
