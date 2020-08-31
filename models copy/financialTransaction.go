package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type FinancialTransaction struct {
	ID                nulls.Int     `col:"" json:"id"`
	TransdateAt       nulls.Time    `col:"" json:"transdateAt"`
	Amount_out        nulls.Float32 `col:"" json:"amount_out"`
	Amount_in         nulls.Float32 `col:"" json:"amount_in"`
	Balance           nulls.Float32 `col:"" json:"balance"`
	Tt_transUse       nulls.String  `col:"" json:"tt_transUse"`
	IsContractPayment nulls.Bool    `col:"" json:"isContractPayment"`
	BankaccountName   nulls.String  `col:"" json:"bankaccountName"`
	BankaccountNo     nulls.String  `col:"" json:"bankaccountNo"`
	Memo              nulls.String  `col:"" json:"memo"`
	UpdateAt          nulls.Time    `col:"newtime" json:"updateAt"`

	Order_form_id            nulls.Int `col:"fk" json:"order_form_id"`
	PaymentType_id           nulls.Int `col:"fk" json:"paymentType_id"`
	FinancialAccount_id      nulls.Int `col:"fk" json:"financialAccount_id" validate:"required" errm:"必填"`
	FinancialLedgerDebit_id  nulls.Int `col:"fk" json:"financialLedgerDebit_id"`
	FinancialLedgerCredit_id nulls.Int `col:"fk" json:"financialLedgerCredit_id"`
	Currency_id              nulls.Int `col:"fk" json:"currency_id"`
	Company_id               nulls.Int `col:"fk" json:"company_id"`
	UpdateUser_id            nulls.Int `col:"fk" json:"updateUser_id"`

	//========fk
	OrderForm             OrderForm        `ref:"order_form,order_form_id" json:"order_form_id.row" validate:"-"`
	PaymentType           CommonItem       `ref:"common_item,paymentType_id" json:"paymentType_id.row" validate:"-"`
	FinancialAccount      FinancialAccount `ref:"financial_account,financialAccount_id" json:"financialAccount_id.row" validate:"-"`
	FinancialLedgerDebit  FinancialLedger  `ref:"financial_ledger,financialLedgerDebit_id" json:"financialLedgerDebit_id.row" validate:"-"`
	FinancialLedgerCredit FinancialLedger  `ref:"financial_ledger,financialLedgerCredit_id" json:"financialLedgerCredit_id.row" validate:"-"`
	CurrencyItem          CommonItem       `ref:"common_item,currency_id" json:"currency_id.row" validate:"-"`
	Company               Company          `ref:"company,company_id" json:"company_id.row" validate:"-"`
	UpdateUser            User             `ref:"user,updateUser_id" json:"updateUser_id.row" validate:"-"`
}

type FinancialTransactionList struct {
	Items []*FinancialTransaction
}

func (item *FinancialTransaction) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.TransdateAt,
		&item.Amount_out,
		&item.Amount_in,
		&item.Balance,
		&item.Tt_transUse,
		&item.IsContractPayment,
		&item.BankaccountName,
		&item.BankaccountNo,
		&item.Memo,
		&item.UpdateAt,
		&item.Order_form_id,
		&item.PaymentType_id,
		&item.FinancialAccount_id,
		&item.FinancialLedgerDebit_id,
		&item.FinancialLedgerCredit_id,
		&item.Currency_id,
		&item.Company_id,
		&item.UpdateUser_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *FinancialTransaction) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkOrderForm := OrderForm{}
	fkPaymentType := CommonItem{}
	fkFinancialAccount := FinancialAccount{}
	fkFinancialLedgerDebit := FinancialLedger{}
	fkFinancialLedgerCredit := FinancialLedger{}
	fkCurrency := CommonItem{}
	fkCompany := Company{}
	fkUpdateUser := User{}

	columns = item.Receivers()
	columns = append(columns, fkOrderForm.Receivers()...)
	columns = append(columns, fkPaymentType.Receivers()...)
	columns = append(columns, fkFinancialAccount.Receivers()...)
	columns = append(columns, fkFinancialLedgerDebit.Receivers()...)
	columns = append(columns, fkFinancialLedgerCredit.Receivers()...)
	columns = append(columns, fkCurrency.Receivers()...)
	columns = append(columns, fkCompany.Receivers()...)
	columns = append(columns, fkUpdateUser.Receivers()...)

	err := r.Scan(columns...)

	item.OrderForm = fkOrderForm
	item.PaymentType = fkPaymentType
	item.FinancialAccount = fkFinancialAccount
	item.FinancialLedgerDebit = fkFinancialLedgerDebit
	item.FinancialLedgerCredit = fkFinancialLedgerCredit
	item.CurrencyItem = fkCurrency
	item.Company = fkCompany
	item.UpdateUser = fkUpdateUser

	return err
}

func (item *FinancialTransaction) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkOrderForm := OrderForm{}
	fkPaymentType := CommonItem{}
	fkFinancialAccount := FinancialAccount{}
	fkFinancialLedgerDebit := FinancialLedger{}
	fkFinancialLedgerCredit := FinancialLedger{}
	fkCurrency := CommonItem{}
	fkCompany := Company{}
	fkUpdateUser := User{}

	columns = item.Receivers()
	columns = append(columns, fkOrderForm.Receivers()...)
	columns = append(columns, fkPaymentType.Receivers()...)
	columns = append(columns, fkFinancialAccount.Receivers()...)
	columns = append(columns, fkFinancialLedgerDebit.Receivers()...)
	columns = append(columns, fkFinancialLedgerCredit.Receivers()...)
	columns = append(columns, fkCurrency.Receivers()...)
	columns = append(columns, fkCompany.Receivers()...)
	columns = append(columns, fkUpdateUser.Receivers()...)

	err := r.Scan(columns...)

	item.OrderForm = fkOrderForm
	item.PaymentType = fkPaymentType
	item.FinancialAccount = fkFinancialAccount
	item.FinancialLedgerDebit = fkFinancialLedgerDebit
	item.FinancialLedgerCredit = fkFinancialLedgerCredit
	item.CurrencyItem = fkCurrency
	item.Company = fkCompany
	item.UpdateUser = fkUpdateUser

	return err
}

func (list *FinancialTransactionList) ScanRow(r *sql.Rows) error {

	item := new(FinancialTransaction) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
