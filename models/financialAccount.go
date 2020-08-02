package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type FinancialAccount struct {
	ID            nulls.Int     `col:"" json:"id"`
	Name          nulls.String  `col:"" json:"name" validate:"required" errm:"必填"`
	Memo          nulls.String  `col:"" json:"memo"`
	Balance       nulls.Float32 `col:"" json:"balance"`
	OriginBalance nulls.Float32 `col:"" json:"originBalance"`

	AccountType nulls.Int `col:"" json:"accountType"`
	Currency_id nulls.Int `col:"fk" json:"currency_id"`
	//========fk
	CurrencyItem CommonItem `ref:"common_item,currency_id" json:"currency_id.row" validate:"-"`
}

type FinancialAccountList struct {
	Items []*FinancialAccount
}

func (item *FinancialAccount) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Name,
		&item.Memo,
		&item.Balance,
		&item.OriginBalance,
		&item.AccountType,
		&item.Currency_id,
	}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *FinancialAccount) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkCurrency := CommonItem{}

	columns = append(item.Receivers(), fkCurrency.Receivers()...)

	err := r.Scan(columns...)

	item.CurrencyItem = fkCurrency

	return err
}

func (item *FinancialAccount) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkCurrency := CommonItem{}

	columns = append(item.Receivers(), fkCurrency.Receivers()...)

	err := r.Scan(columns...)

	item.CurrencyItem = fkCurrency

	return err
}

func (list *FinancialAccountList) ScanRow(r *sql.Rows) error {

	item := new(FinancialAccount) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
