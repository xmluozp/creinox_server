package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Commodity struct {
	ID            nulls.Int    `col:"" json:"id"`
	Code          nulls.String `col:"" json:"code"`
	Name          nulls.String `col:"" json:"name" validate:"required" errm:"必填"`
	UpdateAt      nulls.Time   `col:"newtime" json:"updateAt"`
	CreateAt      nulls.Time   `col:"default" json:"createAt"`
	Memo          nulls.String `col:"" json:"memo"`
	IsDelete      nulls.Bool   `col:"default" json:"isDelete"`
	Category_id   nulls.Int    `col:"fk" json:"category_id"`
	UpdateUser_id nulls.Int    `col:"fk" json:"updateUser_id"`

	Product_id nulls.Int `json:"product_id"`
	// ProductList   ProductList  `col:"" json:"product.rows"`
	Image_id  nulls.Int `col:"fk" json:"image_id,omitempty"`
	ImageItem Image     `ref:"image,image_id" json:"image_id.row" validate:"-"`

	// commoditySell表
	// SellPrice   nulls.Float32 ` json:"sellPrice"`
	// Currency_id nulls.Int     `col:"fk" json:"currency_id"`

	// commodity list的搜索用。searchTerms only, 不是数据库字段。
	CompanyDomesticCustomer_id nulls.Int `json:"companyDomesticCustomer.id"`
	CompanyOverseasCustomer_id nulls.Int `json:"companyOverseasCustomer.id"`
}

type Commodity_product struct {
	Commodity_id nulls.Int `col:"" json:"commodity_id"`
	Product_id   nulls.Int `col:"" json:"product_id"`
	IsMeta       nulls.Int `col:"" json:"isMeta"`
}

type CommodityList struct {
	Items []*Commodity
}

func (item *Commodity) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Code,
		&item.Name,
		&item.UpdateAt,
		&item.CreateAt,
		&item.Memo,
		&item.IsDelete,
		&item.Category_id,
		&item.UpdateUser_id,
		&item.Product_id,
		&item.Image_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *Commodity) ScanRow(r *sql.Row) error {
	var columns []interface{}

	fkImageItem := Image{}

	columns = append(item.Receivers(), fkImageItem.Receivers()...)

	err := r.Scan(columns...)

	item.ImageItem = fkImageItem.Getter()

	return err
}

func (item *Commodity) ScanRows(r *sql.Rows) error {
	var columns []interface{}

	fkImageItem := Image{}

	columns = append(item.Receivers(), fkImageItem.Receivers()...)

	err := r.Scan(columns...)

	item.ImageItem = fkImageItem.Getter()

	return err
}

func (list *CommodityList) ScanRow(r *sql.Rows) error {

	item := new(Commodity) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
