package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type ProductPurchase struct {
	ID            nulls.Int     `col:"" json:"id"`
	ActiveAt      nulls.Time    `col:"" json:"activeAt"`
	ExpireAt      nulls.Time    `col:"" json:"expireAt"`
	BuyPrice      nulls.Float32 `col:"" json:"buyPrice"`
	IsTax         nulls.Bool    `col:"" json:"isTax"`
	IsComponent   nulls.Bool    `col:"" json:"isComponent"`
	Code          nulls.String  `col:"" json:"code"`
	Spec1         nulls.String  `col:"" json:"spec1"`
	Spec2         nulls.String  `col:"" json:"spec2"`
	Spec3         nulls.String  `col:"" json:"spec3"`
	Thickness     nulls.String  `col:"" json:"thickness"`
	UnitWeight    nulls.Float32 `col:"" json:"unitWeight"`
	NetWeight     nulls.Float32 `col:"" json:"netWeight"`
	GrossWeight   nulls.Float32 `col:"" json:"grossWeight"`
	Moq           nulls.Int     `col:"" json:"moq"`
	PackAmount    nulls.Int     `col:"" json:"packAmount"`
	OuterPackL    nulls.Float32 `col:"" json:"outerPackL"`
	OuterPackW    nulls.Float32 `col:"" json:"outerPackW"`
	OuterPackH    nulls.Float32 `col:"" json:"outerPackH"`
	InnerPackL    nulls.Float32 `col:"" json:"innerPackL"`
	InnerPackW    nulls.Float32 `col:"" json:"innerPackW"`
	InnerPackH    nulls.Float32 `col:"" json:"innerPackH"`
	UpdateAt      nulls.Time    `col:"newtime" json:"updateAt"`
	Product_id    nulls.Int     `col:"fk" json:"product_id" validate:"required" errm:"必填"`
	Company_id    nulls.Int     `col:"fk" json:"company_id" validate:"required" errm:"必填"`
	Currency_id   nulls.Int     `col:"fk" json:"currency_id" validate:"required" errm:"必填"`
	Pack_id       nulls.Int     `col:"fk" json:"pack_id"`
	UnitType_id   nulls.Int     `col:"fk" json:"unitType_id"`
	Polishing_id  nulls.Int     `col:"fk" json:"polishing_id"`
	Texture_id    nulls.Int     `col:"fk" json:"texture_id"`
	UpdateUser_id nulls.Int     `col:"fk" json:"updateUser_id"`

	// 显示在列表里

	CompanyItem   Company    `ref:"company,company_id" json:"company_id.row" validate:"-"`
	CurrencyItem  CommonItem `ref:"common_item,currency_id" json:"currency_id.row" validate:"-"`
	PolishingItem CommonItem `ref:"common_item,polishing_id" json:"polishing_id.row" validate:"-"`
	TextureItem   CommonItem `ref:"common_item,texture_id" json:"texture_id.row" validate:"-"`
	PackItem      CommonItem `ref:"common_item,pack_id" json:"pack_id.row" validate:"-"`
}

type ProductPurchaseList struct {
	Items []*ProductPurchase
}

func (item *ProductPurchase) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.ActiveAt,
		&item.ExpireAt,
		&item.BuyPrice,
		&item.IsTax,
		&item.IsComponent,
		&item.Code,
		&item.Spec1,
		&item.Spec2,
		&item.Spec3,
		&item.Thickness,
		&item.UnitWeight,
		&item.NetWeight,
		&item.GrossWeight,
		&item.Moq,
		&item.PackAmount,
		&item.OuterPackL,
		&item.OuterPackW,
		&item.OuterPackH,
		&item.InnerPackL,
		&item.InnerPackW,
		&item.InnerPackH,
		&item.UpdateAt,
		&item.Product_id,
		&item.Company_id,
		&item.Currency_id,
		&item.Pack_id,
		&item.UnitType_id,
		&item.Polishing_id,
		&item.Texture_id,
		&item.UpdateUser_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *ProductPurchase) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkCompany := Company{}
	fkCurrency := CommonItem{}
	fkPolishing := CommonItem{}
	fkTexture := CommonItem{}
	fkPack := CommonItem{}

	columns = append(item.Receivers(), fkCompany.Receivers()...)
	columns = append(columns, fkCurrency.Receivers()...)
	columns = append(columns, fkPolishing.Receivers()...)
	columns = append(columns, fkTexture.Receivers()...)
	columns = append(columns, fkPack.Receivers()...)

	err := r.Scan(columns...)

	item.CompanyItem = fkCompany
	item.CurrencyItem = fkCurrency
	item.PolishingItem = fkPolishing
	item.TextureItem = fkTexture
	item.PackItem = fkPack

	return err
}

func (item *ProductPurchase) ScanRows(r *sql.Rows) error {
	var columns []interface{}

	fkCompany := Company{}
	fkCurrency := CommonItem{}
	fkPolishing := CommonItem{}
	fkTexture := CommonItem{}
	fkPack := CommonItem{}

	columns = append(item.Receivers(), fkCompany.Receivers()...)
	columns = append(columns, fkCurrency.Receivers()...)
	columns = append(columns, fkPolishing.Receivers()...)
	columns = append(columns, fkTexture.Receivers()...)
	columns = append(columns, fkPack.Receivers()...)

	err := r.Scan(columns...)

	item.CompanyItem = fkCompany
	item.CurrencyItem = fkCurrency
	item.PolishingItem = fkPolishing
	item.TextureItem = fkTexture
	item.PackItem = fkPack

	return err
}

func (list *ProductPurchaseList) ScanRow(r *sql.Rows) error {

	item := new(ProductPurchase) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
