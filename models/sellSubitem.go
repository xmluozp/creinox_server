package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type SellSubitem struct {
	ID         nulls.Int     `col:"" json:"id"`
	BuyerCode  nulls.String  `col:"" json:"buyerCode" validate:"required" errm:"必填"`
	BarCode    nulls.String  `col:"" json:"barCode"`
	Amount     nulls.Int     `col:"" json:"amount"`
	PackAmount nulls.Int     `col:"" json:"packAmount"`
	UnitPrice  nulls.Float32 `col:"" json:"unitPrice"`
	Spec       nulls.String  `col:"" json:"spec"`
	Thickness  nulls.String  `col:"" json:"thickness"`

	OuterPackL nulls.Float32 `col:"" json:"outerPackL"`
	OuterPackW nulls.Float32 `col:"" json:"outerPackW"`
	OuterPackH nulls.Float32 `col:"" json:"outerPackH"`
	InnerPackL nulls.Float32 `col:"" json:"innerPackL"`
	InnerPackW nulls.Float32 `col:"" json:"innerPackW"`
	InnerPackH nulls.Float32 `col:"" json:"innerPackH"`

	UnitWeight  nulls.Float32 `col:"" json:"unitWeight"`
	NetWeight   nulls.Float32 `col:"" json:"netWeight"`
	GrossWeight nulls.Float32 `col:"" json:"grossWeight"`

	Fcl20 nulls.Float32 `col:"" json:"fcl20"`
	Fcl40 nulls.Float32 `col:"" json:"fcl40"`

	Commodity_id     nulls.Int `col:"fk" json:"commodity_id" validate:"required" errm:"必填"`
	Sell_contract_id nulls.Int `col:"fk" json:"sell_contract_id" validate:"required" errm:"必填"`
	UnitType_id      nulls.Int `col:"fk" json:"unitType_id"`
	Currency_id      nulls.Int `col:"fk" json:"currency_id" validate:"required" errm:"必填"`
	Polishing_id     nulls.Int `col:"fk" json:"polishing_id"`
	Texture_id       nulls.Int `col:"fk" json:"texture_id"`
	Pack_id          nulls.Int `col:"fk" json:"pack_id"`
	ImagePacking_id  nulls.Int `col:"fk" json:"imagePacking_id"`

	// 显示在列表里
	Commodity    Commodity    `ref:"commodity,commodity_id" json:"commodity_id.row" validate:"-"`
	SellContract SellContract `ref:"combine_sell_contract,sell_contract_id" json:"sell_contract_id.row" validate:"-"` //这里是取combine
	UnitTypeItem CommonItem   `ref:"common_item,unitType_id" json:"unitType_id.row" validate:"-"`
	ImagePacking Image        `ref:"image,imagePacking_id" json:"imagePacking_id.row" validate:"-"`
}

type SellSubitemList struct {
	Items []*SellSubitem
}

func (item *SellSubitem) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.BuyerCode,
		&item.BarCode,
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
		&item.Commodity_id,
		&item.Sell_contract_id,
		&item.UnitType_id,
		&item.Currency_id,
		&item.Polishing_id,
		&item.Texture_id,
		&item.Pack_id,
		&item.ImagePacking_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// 显示view

func (item *SellSubitem) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkCommodity := Commodity{}
	fkSellContract := SellContract{}
	fkUnitTypeItem := CommonItem{}
	fkImagePacking := Image{}

	columns = append(item.Receivers(), fkCommodity.ReceiversOriginal()...)
	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUnitTypeItem.Receivers()...)
	columns = append(columns, fkImagePacking.Receivers()...)

	err := r.Scan(columns...)

	item.Commodity = fkCommodity
	item.SellContract = fkSellContract
	item.UnitTypeItem = fkUnitTypeItem
	item.ImagePacking = fkImagePacking.Getter()

	return err
}

func (item *SellSubitem) ScanRows(r *sql.Rows) error {
	var columns []interface{}

	fkCommodity := Commodity{}
	fkSellContract := SellContract{}
	fkUnitTypeItem := CommonItem{}
	fkImagePacking := Image{}

	columns = append(item.Receivers(), fkCommodity.ReceiversOriginal()...)
	columns = append(columns, fkSellContract.Receivers()...)
	columns = append(columns, fkUnitTypeItem.Receivers()...)
	columns = append(columns, fkImagePacking.Receivers()...)

	err := r.Scan(columns...)

	item.Commodity = fkCommodity
	item.SellContract = fkSellContract
	item.UnitTypeItem = fkUnitTypeItem
	item.ImagePacking = fkImagePacking.Getter()

	return err
}

func (list *SellSubitemList) ScanRow(r *sql.Rows) error {

	item := new(SellSubitem) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
