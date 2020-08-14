package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type MouldContract struct {
	ID   nulls.Int    `col:"" json:"id"`
	Code nulls.String `col:"" json:"code" validate:"required" errm:"必填"`

	Spec             nulls.String  `col:"" json:"spec"`
	UnitPrice        nulls.Float32 `col:"" json:"unitPrice"`
	PrepayPercentage nulls.Float32 `col:"" json:"prepayPercentage"`
	PrepayPrice      nulls.Float32 `col:"" json:"prepayPrice"`

	PrepayAt       nulls.Time `col:"" json:"prepayAt"`
	ActiveAt       nulls.Time `col:"" json:"activeAt"`
	ScheduleAt     nulls.Time `col:"" json:"scheduleAt"`
	DeliverAt      nulls.Time `col:"" json:"deliverAt"`
	Buyer_signAt   nulls.Time `col:"" json:"buyer_signAt"`
	Seller_signAt  nulls.Time `col:"" json:"seller_signAt"`
	DeliverDueDays nulls.Int  `col:"" json:"deliverDueDays"`
	ConfirmDueDays nulls.Int  `col:"" json:"confirmDueDays"`

	Buyer_signer      nulls.String `col:"" json:"buyer_signer"`
	Seller_signer     nulls.String `col:"" json:"seller_signer"`
	Buyer_accountName nulls.String `col:"" json:"buyer_accountName"`
	Buyer_accountNo   nulls.String `col:"" json:"buyer_accountNo"`
	Buyer_bankName    nulls.String `col:"" json:"buyer_bankName"`

	Seller_accountName nulls.String `col:"" json:"seller_accountName"`
	Seller_accountNo   nulls.String `col:"" json:"seller_accountNo"`
	Seller_bankName    nulls.String `col:"" json:"seller_bankName"`
	Tt_memo            nulls.String `col:"" json:"tt_memo"`
	UpdateAt           nulls.Time   `col:"" json:"updateAt"`

	Order_form_id nulls.Int `col:"fk" json:"order_form_id"`
	Product_id    nulls.Int `col:"fk" json:"product_id"`
	Region_id     nulls.Int `col:"fk" json:"region_id"`
	Currency_id   nulls.Int `col:"fk" json:"currency_id"`

	Follower_id       nulls.Int `col:"fk" json:"follower_id"`
	UpdateUser_id     nulls.Int `col:"fk" json:"updateUser_id"`
	Gallary_folder_id nulls.Int `col:"fk" json:"gallary_folder_id"`

	// order里取
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

	// Views 其他地方都是用row，这里实验性的用view
	View_follower            nulls.String `col:"" json:"view_follower"`
	View_productCode         nulls.String `col:"" json:"view_productCode"`
	View_image_thumbnail     nulls.String `col:"" json:"view_image_thumbnail"`
	View_buyer_company_name  nulls.String `col:"" json:"view_buyer_company_name"`
	View_seller_company_name nulls.String `col:"" json:"view_seller_company_name"`

	// 打印用
	Region Region `ref:"region,region_id" json:"region_id.row" validate:"-"`

	// collapse之后要进出款项用
	// BuyContractList []BuyContract `json:"buyContract_list"`
}

type MouldContracList struct {
	Items []*MouldContract
}

func (item *MouldContract) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Code,
		&item.Spec,
		&item.UnitPrice,
		&item.PrepayPercentage,
		&item.PrepayPrice,
		&item.PrepayAt,
		&item.ActiveAt,
		&item.ScheduleAt,
		&item.DeliverAt,
		&item.Buyer_signAt,
		&item.Seller_signAt,
		&item.DeliverDueDays,
		&item.ConfirmDueDays,
		&item.Buyer_signer,
		&item.Seller_signer,
		&item.Buyer_accountName,
		&item.Buyer_accountNo,
		&item.Buyer_bankName,
		&item.Seller_accountName,
		&item.Seller_accountNo,
		&item.Seller_bankName,
		&item.Tt_memo,
		&item.UpdateAt,
		&item.Order_form_id,
		&item.Product_id,
		&item.Region_id,
		&item.Currency_id,
		&item.Follower_id,
		&item.UpdateUser_id,
		&item.Gallary_folder_id,
		&item.ContractType,
		&item.InvoiceCode,
		&item.TotalPrice,
		&item.PaidPrice,
		&item.Seller_company_id,
		&item.Buyer_company_id,
		&item.SellerAddress,
		&item.BuyerAddress,
		&item.IsDone,
		&item.Order_memo,
		// views
		&item.View_follower,
		&item.View_productCode,
		&item.View_image_thumbnail,
		&item.View_buyer_company_name,
		&item.View_seller_company_name}

	// 这里view直接取，因为要做在combine的view里面
	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *MouldContract) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkRegion := Region{}
	columns = item.Receivers()
	columns = append(columns, fkRegion.Receivers()...)
	err := r.Scan(columns...)
	item.Region = fkRegion
	return err
}

func (item *MouldContract) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkRegion := Region{}
	columns = item.Receivers()
	columns = append(columns, fkRegion.Receivers()...)
	err := r.Scan(columns...)
	item.Region = fkRegion

	return err
}

func (list *MouldContracList) ScanRow(r *sql.Rows) error {

	item := new(MouldContract) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
