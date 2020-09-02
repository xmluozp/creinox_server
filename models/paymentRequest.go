package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type PaymentRequest struct {
	ID             nulls.Int     `col:"" json:"id"`
	Amount         nulls.Float32 `col:"" json:"amount"`
	RequestType    nulls.Int     `col:"" json:"requestType"` // 0:合同付款 1: 其他付款
	InvoiceCode    nulls.String  `col:"" json:"invoiceCode"`
	To_companyName nulls.String  `col:"" json:"to_companyName"`

	BankaccountName nulls.String `col:"" json:"bankaccountName"`
	BankaccountNo   nulls.String `col:"" json:"bankaccountNo"`
	Tt_transUse     nulls.String `col:"" json:"tt_transUse"`
	Location        nulls.String `col:"" json:"location"`
	CreateAt        nulls.Time   `col:"" json:"createAt"`
	ExpressAt       nulls.Time   `col:"" json:"expressAt"`
	ExpiryAt        nulls.Time   `col:"" json:"expiryAt"`  // 付款到期日
	ApproveAt       nulls.Time   `col:"" json:"approveAt"` // 审批时间
	Status          nulls.Int    `col:"" json:"status"`    // 0, 申请， 1，通过， 2，拒绝
	Memo            nulls.String `col:"" json:"memo"`

	//========fk
	Order_form_id    nulls.Int `col:"fk" json:"order_form_id"`
	ApplicantUser_id nulls.Int `col:"fk" json:"applicantUser_id"`
	ApproveUser_id   nulls.Int `col:"fk" json:"approveUser_id"`
	From_company_id  nulls.Int `col:"fk" json:"from_company_id"` // 景诚钰诚
	To_company_id    nulls.Int `col:"fk" json:"to_company_id"`   // 选了以后可以选银行，最终只用到文字
	PaymentType_id   nulls.Int `col:"fk" json:"paymentType_id"`
	Currency_id      nulls.Int `col:"fk" json:"currency_id"`

	//========items from fk
	OrderForm       OrderForm  `ref:"order_form,order_form_id" json:"order_form_id.row" validate:"-"`
	ApplicantUser   User       `ref:"user,applicantUser_id" json:"applicantUser_id.row" validate:"-"`
	ApproveUser     User       `ref:"user,approveUser_id" json:"approveUser_id.row" validate:"-"`
	FromCompany     Company    `ref:"company,from_company_id" json:"from_company_id.row" validate:"-"`
	ToCompany       Company    `ref:"company,to_company_id" json:"to_company_id.row" validate:"-"`
	PaymentTypeItem CommonItem `ref:"common_item,paymentType_id" json:"paymentType_id.row" validate:"-"`
	CurrencyItem    CommonItem `ref:"common_item,currency_id" json:"currency_id.row" validate:"-"`
}

func (item *PaymentRequest) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Amount,
		&item.RequestType,
		&item.InvoiceCode,
		&item.To_companyName,
		&item.BankaccountName,
		&item.BankaccountNo,
		&item.Tt_transUse,
		&item.Location,
		&item.CreateAt,
		&item.ExpressAt,
		&item.ExpiryAt,
		&item.ApproveAt,
		&item.Status,
		&item.Memo,
		&item.Order_form_id,
		&item.ApplicantUser_id,
		&item.ApproveUser_id,
		&item.From_company_id,
		&item.To_company_id,
		&item.PaymentType_id,
		&item.Currency_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *PaymentRequest) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkOrderForm := OrderForm{}
	fkApplicant := User{}
	fkApprove := User{}
	fkFromCompany := Company{}
	fkToCompany := Company{}
	fkPaymentTypeItem := CommonItem{}
	fkCurrencyItem := CommonItem{}

	columns = append(item.Receivers(), fkOrderForm.Receivers()...)
	columns = append(columns, fkApplicant.Receivers()...)
	columns = append(columns, fkApprove.Receivers()...)
	columns = append(columns, fkFromCompany.Receivers()...)
	columns = append(columns, fkToCompany.Receivers()...)
	columns = append(columns, fkPaymentTypeItem.Receivers()...)
	columns = append(columns, fkCurrencyItem.Receivers()...)

	err := r.Scan(columns...)

	item.OrderForm = fkOrderForm
	item.ApplicantUser = fkApplicant
	item.ApproveUser = fkApprove
	item.FromCompany = fkFromCompany
	item.ToCompany = fkToCompany
	item.PaymentTypeItem = fkPaymentTypeItem
	item.CurrencyItem = fkCurrencyItem

	return err
}

func (item *PaymentRequest) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkOrderForm := OrderForm{}
	fkApplicant := User{}
	fkApprove := User{}
	fkFromCompany := Company{}
	fkToCompany := Company{}
	fkPaymentTypeItem := CommonItem{}
	fkCurrencyItem := CommonItem{}

	columns = append(item.Receivers(), fkOrderForm.Receivers()...)
	columns = append(columns, fkApplicant.Receivers()...)
	columns = append(columns, fkApprove.Receivers()...)
	columns = append(columns, fkFromCompany.Receivers()...)
	columns = append(columns, fkToCompany.Receivers()...)
	columns = append(columns, fkPaymentTypeItem.Receivers()...)
	columns = append(columns, fkCurrencyItem.Receivers()...)

	err := r.Scan(columns...)

	item.OrderForm = fkOrderForm
	item.ApplicantUser = fkApplicant
	item.ApproveUser = fkApprove
	item.FromCompany = fkFromCompany
	item.ToCompany = fkToCompany
	item.PaymentTypeItem = fkPaymentTypeItem
	item.CurrencyItem = fkCurrencyItem

	return err
}
