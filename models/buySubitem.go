package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type BuySubitem struct {
	ID         nulls.Int    `col:"" json:"id"`
	SellerCode nulls.String `col:"" json:"sellerCode" validate:"required" errm:"必填"`

	IsReceipt  nulls.Bool    `col:"" json:"isReceipt"`
	Amount     nulls.Int     `col:"" json:"amount"`
	PackAmount nulls.Int     `col:"" json:"packAmount"`
	UnitPrice  nulls.Float32 `col:"" json:"unitPrice"`

	Spec        nulls.String  `col:"" json:"spec"`
	Thickness   nulls.Float32 `col:"" json:"thickness"`
	UnitWeight  nulls.Float32 `col:"" json:"unitWeight"`
	NetWeight   nulls.Float32 `col:"" json:"netWeight"`
	GrossWeight nulls.Float32 `col:"" json:"grossWeight"`

	OuterPackL   nulls.Float32 `col:"" json:"outerPackL"`
	OuterPackW   nulls.Float32 `col:"" json:"outerPackW"`
	OuterPackH   nulls.Float32 `col:"" json:"outerPackH"`
	InnerPackL   nulls.Float32 `col:"" json:"innerPackL"`
	InnerPackW   nulls.Float32 `col:"" json:"innerPackW"`
	InnerPackH   nulls.Float32 `col:"" json:"innerPackH"`
	Fcl20        nulls.Float32 `col:"" json:"fcl20"`
	Fcl40        nulls.Float32 `col:"" json:"fcl40"`
	PickuptimeAt nulls.Time    `col:"" json:"pickuptimeAt"`

	Product_id      nulls.Int `col:"fk" json:"product_id" validate:"required" errm:"必填"`
	Sell_subitem_id nulls.Int `col:"fk" json:"sell_subitem_id"`
	Buy_contract_id nulls.Int `col:"fk" json:"buy_contract_id" validate:"required" errm:"必填"`
	UnitType_id     nulls.Int `col:"fk" json:"unitType_id"`
	Currency_id     nulls.Int `col:"fk" json:"currency_id" validate:"required" errm:"必填"`
	Polishing_id    nulls.Int `col:"fk" json:"polishing_id"`
	Texture_id      nulls.Int `col:"fk" json:"texture_id"`
	Pack_id         nulls.Int `col:"fk" json:"pack_id"`

	// 显示在列表里
	Product      Product     `ref:"product,product_id" json:"product_id.row" validate:"-"`
	BuyContract  BuyContract `ref:"combine_buy_contract,buy_contract_id" json:"buy_contract_id.row" validate:"-"`
	UnitTypeItem CommonItem  `ref:"common_item,unitType_id" json:"unitType_id.row" validate:"-"`
	SellSubitem  SellSubitem `ref:"sell_subitem,sell_subitem_id" json:"sell_subitem_id.row" validate:"-"`
}

type BuySubitemList struct {
	Items []*BuySubitem
}

func (item *BuySubitem) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.SellerCode,
		&item.IsReceipt,
		&item.Amount,
		&item.PackAmount,
		&item.UnitPrice,

		&item.Spec,
		&item.Thickness,
		&item.UnitWeight,
		&item.NetWeight,
		&item.GrossWeight,

		&item.OuterPackL,
		&item.OuterPackW,
		&item.OuterPackH,
		&item.InnerPackL,
		&item.InnerPackW,
		&item.InnerPackH,
		&item.Fcl20,
		&item.Fcl40,
		&item.PickuptimeAt,

		&item.Product_id,
		&item.Sell_subitem_id,
		&item.Buy_contract_id,
		&item.UnitType_id,
		&item.Currency_id,
		&item.Polishing_id,
		&item.Texture_id,
		&item.Pack_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// 显示view

func (item *BuySubitem) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkProduct := Product{}
	fkBuyContract := BuyContract{}
	fkUnitTypeItem := CommonItem{}
	fkSellSubitem := SellSubitem{}

	columns = append(item.Receivers(), fkProduct.Receivers()...)
	columns = append(columns, fkBuyContract.Receivers()...)
	columns = append(columns, fkUnitTypeItem.Receivers()...)
	columns = append(columns, fkSellSubitem.Receivers()...)

	err := r.Scan(columns...)

	item.Product = fkProduct
	item.BuyContract = fkBuyContract
	item.UnitTypeItem = fkUnitTypeItem
	item.SellSubitem = fkSellSubitem

	return err
}

func (item *BuySubitem) ScanRows(r *sql.Rows) error {
	var columns []interface{}

	fkProduct := Product{}
	fkBuyContract := BuyContract{}
	fkUnitTypeItem := CommonItem{}
	fkSellSubitem := SellSubitem{}

	columns = append(item.Receivers(), fkProduct.Receivers()...)
	columns = append(columns, fkBuyContract.Receivers()...)
	columns = append(columns, fkUnitTypeItem.Receivers()...)
	columns = append(columns, fkSellSubitem.Receivers()...)

	err := r.Scan(columns...)

	item.Product = fkProduct
	item.BuyContract = fkBuyContract
	item.UnitTypeItem = fkUnitTypeItem
	item.SellSubitem = fkSellSubitem

	return err
}

func (list *BuySubitemList) ScanRow(r *sql.Rows) error {

	item := new(BuySubitem) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
