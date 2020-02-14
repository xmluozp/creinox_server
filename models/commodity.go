package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Commodity struct {
	ID            nulls.Int    `col:"" json:"id"`
	Name          nulls.String `col:"" json:"name" validate:"required" errm:"必填"`
	Memo          nulls.String `col:"" json:"memo"`
	UpdateAt      nulls.Time   `col:"newtime" json:"updateAt"`
	CreateAt      nulls.Time   `col:"default" json:"createAt"`
	IsDelete      nulls.Bool   `col:"default" json:"isDelete"`
	UpdateUser_id nulls.Int    `col:"fk" json:"updateUser_id"`
	Code          nulls.String `col:"" json:"code"`
	Category_id   nulls.Int    `col:"fk" json:"category_id"`
	Product_id    nulls.Int    `col:"fk" json:"product_id"`
	ProductList   ProductList  `col:"" json:"product.rows"`
	Image_id      nulls.Int    `col:"fk" json:"image_id"`
	ImageItem     Image        `col:"" json:"image_id.row"`

	// commoditySell表
	SellPrice   nulls.Float32 `col:"" json:"sellPrice"`
	Currency_id nulls.Int     `col:"fk" json:"currency_id"`

	// commodity list的搜索用。searchTerms only, 不是数据库字段。
	CompanyDomesticCustomer_id nulls.Int `col:"fk" json:"companyDomesticCustomer.id"`
	CompanyOverseasCustomer_id nulls.Int `col:"fk" json:"companyOverseasCustomer.id"`
}

type CommodityList struct {
	Items []*Commodity
}

func (item *Commodity) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Name,
		&item.Memo,
		&item.UpdateAt,
		&item.CreateAt,
		&item.IsDelete,
		&item.UpdateUser_id,
		&item.Code,
		&item.Category_id,
		&item.Product_id,
		&item.ProductList,
		&item.Image_id,
		&item.ImageItem,
		&item.SellPrice,
		&item.Currency_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *Commodity) ScanRow(r *sql.Row) error {
	return r.Scan(item.Receivers()...)
}

func (item *Commodity) ScanRows(r *sql.Rows) error {
	return r.Scan(item.Receivers()...)
}

func (list *CommodityList) ScanRow(r *sql.Rows) error {

	item := new(Commodity) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
