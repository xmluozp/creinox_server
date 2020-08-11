package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

// 这个model单纯为了保存，不负责显示
type OrderForm struct {
	ID             nulls.Int     `col:"" json:"id"`
	ContractType   nulls.Int     `col:"" json:"contractType,omitempty"`
	Code           nulls.String  `col:"" json:"code"`
	InvoiceCode    nulls.String  `col:"" json:"invoiceCode"`
	Receivable     nulls.Float32 `col:"" json:"receivable"`
	ReceivablePaid nulls.Float32 `col:"" json:"receivablePaid"`
	Payable        nulls.Float32 `col:"" json:"payable"`
	PayablePaid    nulls.Float32 `col:"" json:"payablePaid"`

	// 搜索用
	KeyWord nulls.String `json:"keyword" keywords:"code|order_memo|invoiceCode"`

	Seller_company_id nulls.Int    `col:"fk" json:"seller_company_id"`
	Buyer_company_id  nulls.Int    `col:"fk" json:"buyer_company_id"`
	IsDone            nulls.Bool   `col:"" json:"isDone"`
	Order_memo        nulls.String `col:"" json:"order_memo"`
}
type OrderFormList struct {
	Items []*OrderForm
}

func (item *OrderForm) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.ContractType,
		&item.Code,
		&item.InvoiceCode,
		&item.Receivable,
		&item.ReceivablePaid,
		&item.Payable,
		&item.PayablePaid,
		&item.Seller_company_id,
		&item.Buyer_company_id,
		&item.IsDone,
		&item.Order_memo}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *OrderForm) ScanRow(r *sql.Row) error {

	err := r.Scan(item.Receivers()...)
	return err
}

func (item *OrderForm) ScanRows(r *sql.Rows) error {

	err := r.Scan(item.Receivers()...)
	return err
}

func (list *OrderFormList) ScanRow(r *sql.Rows) error {

	item := new(OrderForm) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
