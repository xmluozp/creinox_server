package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Product struct {
	ID            nulls.Int     `col:"" json:"id"`
	Code          nulls.String  `col:"" json:"code" validate:"required" errm:"必填"`
	Name          nulls.String  `col:"" json:"name"`
	EName         nulls.String  `col:"" json:"ename"`
	Shortname     nulls.String  `col:"" json:"shortname"`
	EShortname    nulls.String  `col:"" json:"eshortname"`
	Spec1         nulls.String  `col:"" json:"spec1"`
	Spec2         nulls.String  `col:"" json:"spec2"`
	Spec3         nulls.String  `col:"" json:"spec3"`
	Barcode       nulls.String  `col:"" json:"barcode"`
	Thickness     nulls.Float32 `col:"" json:"thickness"`
	UnitWeight    nulls.Float32 `col:"" json:"unitWeight"`
	RetrieveTime  nulls.Time    `col:"" json:"retrieveTime"`
	UpdateAt      nulls.Time    `col:"newtime" json:"updateAt"`
	CreateAt      nulls.Time    `col:"default" json:"createAt"`
	Memo          nulls.String  `col:"" json:"memo"`
	IsOEM         nulls.Bool    `col:"" json:"isOEM"`
	IsSemiProduct nulls.Bool    `col:"" json:"isSemiProduct"`
	IsEndProduct  nulls.Bool    `col:"" json:"isEndProduct"`
	IsDelete      nulls.Bool    `col:"default" json:"isDelete"`
	Polishing_id  nulls.Int     `col:"fk" json:"polishing_id"`
	Texture_id    nulls.Int     `col:"fk" json:"texture_id"`
	Retriever_id  nulls.Int     `col:"fk" json:"retriever_id"`
	UpdateUser_id nulls.Int     `col:"fk" json:"updateUser_id"`
	Category_id   nulls.Int     `col:"fk" json:"category_id"`
	Image_id      nulls.Int     `col:"fk" json:"image_id,omitempty"`
	BuyPrice      nulls.Float32 `json:"buyPrice"`
	Currency_id   nulls.Int     `json:"currency_id"`

	ComodityCode      nulls.String `json:"comodity.code"`
	CompanyFactoryId  nulls.Int    `json:"companyFactory.id"`
	IsCreateCommodity nulls.Bool   `json:"isCreateCommodity"`

	ImageItem Image `ref:"image,image_id" json:"image_id.row" validate:"-"`
}
type ProductList struct {
	Items []*Product
}

func (item *Product) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Code,
		&item.Name,
		&item.EName,
		&item.Shortname,
		&item.EShortname,
		&item.Spec1,
		&item.Spec2,
		&item.Spec3,
		&item.Barcode,
		&item.Thickness,
		&item.UnitWeight,
		&item.RetrieveTime,
		&item.UpdateAt,
		&item.CreateAt,
		&item.Memo,
		&item.IsOEM,
		&item.IsSemiProduct,
		&item.IsEndProduct,
		&item.IsDelete,
		&item.Polishing_id,
		&item.Texture_id,
		&item.Retriever_id,
		&item.UpdateUser_id,
		&item.Category_id,
		&item.Image_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *Product) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkImageItem := Image{}

	columns = append(item.Receivers(), fkImageItem.Receivers()...)

	err := r.Scan(columns...)

	item.ImageItem = fkImageItem.Getter()

	return err
}

func (item *Product) ScanRows(r *sql.Rows) error {
	var columns []interface{}

	fkImageItem := Image{}

	columns = append(item.Receivers(), fkImageItem.Receivers()...)

	err := r.Scan(columns...)

	item.ImageItem = fkImageItem.Getter()

	return err
}

func (list *ProductList) ScanRow(r *sql.Rows) error {

	item := new(Product) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
