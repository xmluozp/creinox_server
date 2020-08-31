package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type ExpressOrder struct {
	ID                  nulls.Int    `col:"" json:"id"`
	Code                nulls.String `col:"" json:"code" validate:"required" errm:"必填"`
	Direction           nulls.Int    `col:"" json:"direction"`   // 0 寄件 1 收件
	ExpressType         nulls.Int    `col:"" json:"expressType"` // 0文件，1包裹
	CreateAt            nulls.Time   `col:"" json:"createAt"`
	ExpressAt           nulls.Time   `col:"" json:"expressAt"`
	Memo                nulls.String `col:"" json:"memo"`
	Internal_company_id nulls.Int    `col:"fk" json:"internal_company_id"`
	External_company_id nulls.Int    `col:"fk" json:"external_company_id"`
	ExpressCompany_id   nulls.Int    `col:"fk" json:"expressCompany_id"`

	//========fk
	InternalCompany Company    `ref:"company,internal_company_id" json:"internal_company_id.row" validate:"-"`
	ExternalCompany Company    `ref:"company,external_company_id" json:"external_company_id.row" validate:"-"`
	ExpressCompany  CommonItem `ref:"common_item,expressCompany_id" json:"expressCompany_id.row" validate:"-"`
}

func (item *ExpressOrder) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Code,
		&item.Direction,
		&item.ExpressType,
		&item.CreateAt,
		&item.ExpressAt,
		&item.Memo,
		&item.Internal_company_id,
		&item.External_company_id,
		&item.ExpressCompany_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *ExpressOrder) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkInternalCompany := Company{}
	fkExternalCompany := Company{}
	fkExpressCompany := CommonItem{}

	columns = append(item.Receivers(), fkInternalCompany.Receivers()...)
	columns = append(columns, fkExternalCompany.Receivers()...)
	columns = append(columns, fkExpressCompany.Receivers()...)

	err := r.Scan(columns...)

	item.InternalCompany = fkInternalCompany
	item.ExternalCompany = fkExternalCompany
	item.ExpressCompany = fkExpressCompany

	return err
}

func (item *ExpressOrder) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkInternalCompany := Company{}
	fkExternalCompany := Company{}
	fkExpressCompany := CommonItem{}

	columns = append(item.Receivers(), fkInternalCompany.Receivers()...)
	columns = append(columns, fkExternalCompany.Receivers()...)
	columns = append(columns, fkExpressCompany.Receivers()...)

	err := r.Scan(columns...)

	item.InternalCompany = fkInternalCompany
	item.ExternalCompany = fkExternalCompany
	item.ExpressCompany = fkExpressCompany

	return err
}
