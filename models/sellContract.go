package models

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
)

type SellContract struct {
	ID          nulls.Int    `col:"" json:"id"`
	Code        nulls.String `col:"" json:"code" validate:"required" errm:"必填"`
	OrderNumber nulls.String `col:"" json:"orderNumber"`
	ActiveAt    nulls.Time   `col:"" json:"activeAt"`
	DeliverAt   nulls.Time   `col:"" json:"deliverAt"`
	UpdateAt    nulls.Time   `col:"newtime" json:"updateAt"`

	IsInBaches  nulls.Bool `col:"" json:"isInBaches"`
	IsTransport nulls.Bool `col:"" json:"isTransport"`

	ShippingPrice nulls.Float32 `col:"" json:"shippingPrice"`
	ComissionType nulls.Int     `col:"fk" json:"comissionType"`

	Tt_shipmentDue      nulls.String `col:"" json:"tt_shipmentDue"`
	Tt_insurance        nulls.String `col:"" json:"tt_insurance"`
	Tt_paymentCondition nulls.String `col:"" json:"tt_paymentCondition"`

	UpdateUser_id nulls.Int `col:"fk" json:"updateUser_id"`
	Follower_id   nulls.Int `col:"fk" json:"follower_id"`

	Seller_company_id nulls.Int `col:"fk" json:"seller_company_id" validate:"required" errm:"必填"`
	Buyer_company_id  nulls.Int `col:"fk" json:"buyer_company_id" validate:"required" errm:"必填"`

	ShippingType_id      nulls.Int `col:"fk" json:"shippingType_id"`
	PricingTerm_id       nulls.Int `col:"fk" json:"pricingTerm_id"`
	PaymentType_id       nulls.Int `col:"fk" json:"paymentType_id"`
	Commission_id        nulls.Int `col:"fk" json:"commission_id"`
	Region_id            nulls.Int `col:"fk" json:"region_id"`
	Departure_port_id    nulls.Int `col:"fk" json:"departure_port_id"`
	Destination_port_id  nulls.Int `col:"fk" json:"destination_port_id"`
	Currency_id          nulls.Int `col:"fk" json:"currency_id"`
	Shipping_currency_id nulls.Int `col:"fk" json:"shipping_currency_id"`

	// 显示在列表里
	UserFollower User    `ref:"user,follower_id" json:"follower_id.row" validate:"-"`
	CompanyBuyer Company `ref:"company,buyer_company_id" json:"buyer_company_id.row" validate:"-"`

	View_totalPrice nulls.Float32 `json:"view_totalPrice"`

	// collapse的对应合同列表
	BuyContractList []BuyContract `json:"buyContract_list"`
	SellSubitem     []SellSubitem `json:"subitem_list"`

	// ModelContractList ModelContractList `json:"modelContract_list"`
}

type SellContractList struct {
	Items []*SellContract
}

func (item *SellContract) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Code,
		&item.OrderNumber,
		&item.ActiveAt,
		&item.DeliverAt,
		&item.IsInBaches,
		&item.IsTransport,
		&item.ShippingPrice,
		&item.ComissionType,
		&item.Tt_shipmentDue,
		&item.Tt_insurance,
		&item.Tt_paymentCondition,
		&item.UpdateAt,
		&item.UpdateUser_id,
		&item.Follower_id,
		&item.Seller_company_id,
		&item.Buyer_company_id,
		&item.ShippingType_id,
		&item.PricingTerm_id,
		&item.PaymentType_id,
		&item.Commission_id,
		&item.Region_id,
		&item.Departure_port_id,
		&item.Destination_port_id,
		&item.Currency_id,
		&item.Shipping_currency_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *SellContract) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkUserFollower := User{}
	fkCompanyBuyer := Company{}

	columns = append(item.Receivers(), fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanyBuyer.Receivers()...)

	err := r.Scan(columns...)

	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer

	return err
}

func (item *SellContract) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkUserFollower := User{}
	fkCompanyBuyer := Company{}

	columns = append(item.Receivers(), fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanyBuyer.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条子订单出错", err.Error)
	}

	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer

	return err
}

func (item *SellContract) ScanRowsView(r *sql.Rows) error {

	var columns []interface{}

	fkUserFollower := User{}
	fkCompanyBuyer := Company{}

	columns = append(item.Receivers(), &item.View_totalPrice)
	columns = append(columns, fkUserFollower.Receivers()...)

	columns = append(columns, fkCompanyBuyer.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条view出错", err.Error)
	}

	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer

	return err
}

func (item *SellContract) ScanRowView(r *sql.Row) error {

	var columns []interface{}

	fkUserFollower := User{}
	fkCompanyBuyer := Company{}

	columns = append(item.Receivers(), &item.View_totalPrice)
	columns = append(columns, fkUserFollower.Receivers()...)

	columns = append(columns, fkCompanyBuyer.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条view出错", err.Error)
	}

	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer

	return err
}

func (list *SellContractList) ScanRow(r *sql.Rows) error {

	item := new(SellContract) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}