package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type BankAccount struct {
	ID          nulls.Int    `col:"" json:"id"`
	AccountName nulls.String `col:"" json:"accountName" validate:"required" errm:"必填"`
	AccountNo   nulls.String `col:"" json:"accountNo"`
	BankName    nulls.String `col:"" json:"bankName"`
	Address     nulls.String `col:"" json:"address"`
	SwiftCode   nulls.String `col:"" json:"swiftCode"`
	BankType    nulls.Int    `col:"" json:"bankType"`
	Memo        nulls.String `col:"" json:"memo"`
	Currency_id nulls.Int    `col:"fk" json:"currency_id"`
	Company_id  nulls.Int    `col:"fk" json:"company_id"`
	//========fk
	CurrencyItem CommonItem `ref:"common_item,currency_id" json:"currency_id.row"`
	CompanyItem  Company    `ref:"company,company_id" json:"company_id.row"`
}

type BankAccountList struct {
	Items []*BankAccount
}

func (item *BankAccount) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.AccountName,
		&item.AccountNo,
		&item.BankName,
		&item.Address,
		&item.SwiftCode,
		&item.BankType,
		&item.Memo,
		&item.Currency_id,
		&item.Company_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *BankAccount) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkCurrency := CommonItem{}
	fkCompany := Company{}

	columns = append(item.Receivers(), fkCurrency.Receivers()...)
	columns = append(columns, fkCompany.Receivers()...)

	err := r.Scan(columns...)

	item.CurrencyItem = fkCurrency
	item.CompanyItem = fkCompany

	return err
}

func (item *BankAccount) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkCurrency := CommonItem{}
	fkCompany := Company{}

	columns = append(item.Receivers(), fkCurrency.Receivers()...)
	columns = append(columns, fkCompany.Receivers()...)

	err := r.Scan(columns...)

	item.CurrencyItem = fkCurrency
	item.CompanyItem = fkCompany

	return err
}

func (list *BankAccountList) ScanRow(r *sql.Rows) error {

	item := new(BankAccount) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
