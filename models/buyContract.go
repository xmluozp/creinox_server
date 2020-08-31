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
	Tt_memo                nulls.String `col:"" json:"tt_memo"`
	UpdateAt               nulls.Time   `col:"newtime" json:"updateAt"`

	Region_id        nulls.Int `col:"fk" json:"region_id"`
	Sell_contract_id nulls.Int `col:"fk" json:"sell_contract_id"`
	PaymentType_id   nulls.Int `col:"fk" json:"paymentType_id"`

	Follower_id   nulls.Int `col:"fk" json:"follower_id"`
	UpdateUser_id nulls.Int `col:"fk" json:"updateUser_id"`
	Order_form_id nulls.Int `col:"fk" json:"order_form_id"`

	// order里取 (应付应收直接统一成totalPrice)
	ContractType      nulls.Int     `json:"contractType"`
	InvoiceCode       nulls.String  `json:"invoiceCode"`
	TotalPrice        nulls.Float32 `json:"totalPrice"`
	PaidPrice         nulls.Float32 `json:"paidPrice"`
	Seller_company_id nulls.Int     `json:"seller_company_id"`
	Buyer_company_id  nulls.Int     `json:"buyer_company_id"`
	SellerAddress     nulls.String  `json:"sellerAddress"`
	BuyerAddress      nulls.String  `json:"buyerAddress"`
	IsDone            nulls.Bool    `json:"isDone"`
	Order_memo        nulls.String  `json:"order_memo"`

	// 显示在列表里
	SellContract  SellContract `ref:"combine_sell_contract,sell_contract_id" json:"sell_contract_id.row" validate:"-"`
	UserFollower  User         `ref:"user,follower_id" json:"follower_id.row" validate:"-"`
	CompanyBuyer  Company      `ref:"company,buyer_company_id" json:"buyer_company_id.row" validate:"-"`
	CompanySeller Company      `ref:"company,seller_company_id" json:"seller_company_id.row" validate:"-"`
	Region        Region       `ref:"region,region_id" json:"region_id.row" validate:"-"`

	// collapse的对应子合同列表（放的是产品）
	BuySubitem               []BuySubitem           `json:"buy_subitem_list"`
	FinancialTransactionList []FinancialTransaction `json:"financialTransaction_list"`
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
		&item.Tt_memo,
		&item.UpdateAt,
		&item.Order_form_id,
		&item.Region_id,
		&item.Sell_contract_id,
		&item.PaymentType_id,
		&item.Follower_id,
		&item.UpdateUser_id,
		&item.ContractType,
		&item.InvoiceCode,
		&item.TotalPrice,
		&item.PaidPrice,
		&item.Seller_company_id,
		&item.Buyer_company_id,
		&item.SellerAddress,
		&item.BuyerAddress,
		&item.IsDone,
		&item.Order_memo}

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
	fkCompanyBuyer := Company{}
	fkCompanySeller := Company{}
	fkRegion := Region{}

	columns = item.Receivers()
	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanyBuyer.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)
	columns = append(columns, fkRegion.Receivers()...)

	err := r.Scan(columns...)

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer
	item.CompanySeller = fkCompanySeller
	item.Region = fkRegion

	return err
}

func (item *BuyContract) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanyBuyer := Company{}
	fkCompanySeller := Company{}
	fkRegion := Region{}

	columns = item.Receivers()
	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanyBuyer.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)
	columns = append(columns, fkRegion.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条view出错", err.Error)
	}

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer
	item.CompanySeller = fkCompanySeller
	item.Region = fkRegion

	return err
}

func (item *BuyContract) ScanRowsView(r *sql.Rows) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanyBuyer := Company{}
	fkCompanySeller := Company{}
	fkRegion := Region{}

	columns = item.Receivers()
	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanyBuyer.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)
	columns = append(columns, fkRegion.Receivers()...)

	err := r.Scan(columns...)

	if err != nil {
		fmt.Println("读取多条view出错", err.Error)
	}

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer
	item.CompanySeller = fkCompanySeller
	item.Region = fkRegion

	return err
}

func (item *BuyContract) ScanRowView(r *sql.Row) error {

	var columns []interface{}

	fkSellContract := SellContract{}
	fkUserFollower := User{}
	fkCompanyBuyer := Company{}
	fkCompanySeller := Company{}
	fkRegion := Region{}

	columns = item.Receivers()
	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUserFollower.Receivers()...)
	columns = append(columns, fkCompanyBuyer.Receivers()...)
	columns = append(columns, fkCompanySeller.Receivers()...)
	columns = append(columns, fkRegion.Receivers()...)

	err := r.Scan(columns...)

	item.SellContract = fkSellContract
	item.UserFollower = fkUserFollower
	item.CompanyBuyer = fkCompanyBuyer
	item.CompanySeller = fkCompanySeller
	item.Region = fkRegion

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
