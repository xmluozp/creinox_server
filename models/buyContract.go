package models

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
)

type BuyContract struct {
	ID        nulls.Int    `col:"" json:"id"`
	Code      nulls.String `col:"" json:"code" validate:"required" errm:"必填"`
	ActiveAt  nulls.Time   `col:"" json:"activeAt"`
	DeliverAt nulls.Time   `col:"" json:"deliverAt"`

	Tt_quality             nulls.String `col:"" json:"tt_quality"`
	Tt_deliveryMethod      nulls.String `col:"" json:"tt_deliveryMethod"`
	Tt_shippingTerm        nulls.String `col:"" json:"tt_shippingTerm"`
	Tt_loss                nulls.String `col:"" json:"tt_loss"`
	Tt_packingStandard     nulls.String `col:"" json:"tt_packingStandard"`
	Tt_acceptanceCondition nulls.String `col:"" json:"tt_acceptanceCondition"`
	Tt_accessories         nulls.String `col:"" json:"tt_accessories"`
	Tt_payment             nulls.String `col:"" json:"tt_payment"`
	Tt_breach              nulls.String `col:"" json:"tt_breach"`
	Tt_dispute             nulls.String `col:"" json:"tt_dispute"`
	Memo                   nulls.String `col:"" json:"memo"`
	UpdateAt               nulls.Time   `col:"newtime" json:"updateAt"`

	Region_id         nulls.Int `col:"fk" json:"region_id"`
	Sell_contract_id  nulls.Int `col:"fk" json:"sell_contract_id"`
	PaymentType_id    nulls.Int `col:"fk" json:"paymentType_id"`
	Seller_company_id nulls.Int `col:"fk" json:"seller_company_id" validate:"required" errm:"必填"`

	Follower_id   nulls.Int `col:"fk" json:"follower_id"`
	UpdateUser_id nulls.Int `col:"fk" json:"updateUser_id"`

	// 显示在列表里
	SellContract  SellContract `ref:"sell_contract,sell_contract_id" json:"sell_contract_id.row" validate:"-"`
	UserFollower  User         `ref:"user,follower_id" json:"follower_id.row" validate:"-"`
	CompanySeller Company      `ref:"company,seller_company_id" json:"seller_company_id.row" validate:"-"`

	View_totalPrice nulls.Float32 `json:"view_totalPrice"`
}

type BuyContractList struct {
	Items []*BuyContract
}

func (item *BuyContract) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Code,
		&item.ActiveAt,
		&item.DeliverAt,
		&item.Tt_quality,
		&item.Tt_deliveryMethod,
		&item.Tt_shippingTerm,
		&item.Tt_loss,
		&item.Tt_packingStandard,
		&item.Tt_acceptanceCondition,
		&item.Tt_accessories,
		&item.Tt_payment,
		&item.Tt_breach,
		&item.Tt_dispute,
		&item.Memo,
		&item.UpdateAt,
		&item.Region_id,
		&item.Sell_contract_id,
		&item.PaymentType_id,
		&item.Seller_company_id,
		&item.Follower_id,
		&item.UpdateUser_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *BuyContract) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanySeller := Company{}

	columns = append(item.Receivers(), fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)

	err := r.Scan(columns...)

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanySeller = fkCompanySeller

	return err
}

func (item *BuyContract) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanySeller := Company{}

	columns = append(item.Receivers(), fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条子订单出错", err.Error)
	}

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanySeller = fkCompanySeller

	return err
}

func (item *BuyContract) ScanRowsView(r *sql.Rows) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanySeller := Company{}

	columns = append(item.Receivers(), &item.View_totalPrice)

	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条view出错", err.Error)
	}

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanySeller = fkCompanySeller

	return err
}

func (item *BuyContract) ScanRowView(r *sql.Row) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanySeller := Company{}

	columns = append(item.Receivers(), &item.View_totalPrice)

	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条view出错", err.Error)
	}

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanySeller = fkCompanySeller

	return err
}

func (list *BuyContractList) ScanRow(r *sql.Rows) error {

	item := new(BuyContract) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
