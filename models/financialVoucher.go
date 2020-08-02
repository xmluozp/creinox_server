package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type FinancialVoucher struct {
	ID               nulls.Int     `col:"" json:"id"`
	Resource_code    nulls.String  `col:"" json:"resource_code"`
	Word             nulls.String  `col:"" json:"word"`
	Number           nulls.String  `col:"" json:"number"`
	Debit            nulls.Float32 `col:"" json:"debit"`
	Credit           nulls.Float32 `col:"" json:"credit"`
	FinancialSubject nulls.String  `col:"" json:"financialSubject"`
	DetailedSubject  nulls.String  `col:"" json:"detailedSubject"`
	Memo             nulls.String  `col:"" json:"memo"`
	CreateAt         nulls.Time    `col:"default" json:"createAt"`

	FinancialAccount_id nulls.Int `col:"fk" json:"financialAccount_id"`
	UpdateUser_id       nulls.Int `col:"fk" json:"updateUser_id"`

	//========fk
	// 币种是根据financial Account决定的
	FinancialAccount FinancialAccount `ref:"financial_account,financialAccount_id" json:"financialAccount_id.row" validate:"-"`
	UpdateUser       User             `ref:"user,updateUser_id" json:"updateUser_id.row" validate:"-"`
}

type FinancialVoucherList struct {
	Items []*FinancialVoucher
}

func (item *FinancialVoucher) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Resource_code,
		&item.Word,
		&item.Number,
		&item.Debit,
		&item.Credit,
		&item.FinancialSubject,
		&item.DetailedSubject,
		&item.Memo,
		&item.CreateAt,
		&item.FinancialAccount_id,
		&item.UpdateUser_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *FinancialVoucher) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkFinancialAccount := FinancialAccount{}
	fkUpdateUser := User{}

	columns = item.Receivers()

	columns = append(columns, fkFinancialAccount.Receivers()...)
	columns = append(columns, fkUpdateUser.Receivers()...)

	err := r.Scan(columns...)

	item.FinancialAccount = fkFinancialAccount
	item.UpdateUser = fkUpdateUser

	return err
}

func (item *FinancialVoucher) ScanRows(r *sql.Rows) error {

	var columns []interface{}
	fkFinancialAccount := FinancialAccount{}
	fkUpdateUser := User{}
	columns = item.Receivers()
	columns = append(columns, fkFinancialAccount.Receivers()...)
	columns = append(columns, fkUpdateUser.Receivers()...)

	err := r.Scan(columns...)

	item.FinancialAccount = fkFinancialAccount
	item.UpdateUser = fkUpdateUser

	return err
}

func (list *FinancialVoucherList) ScanRow(r *sql.Rows) error {

	item := new(FinancialVoucher) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
