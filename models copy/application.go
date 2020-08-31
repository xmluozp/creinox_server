package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Application struct {
	ID               nulls.Int    `col:"" json:"id"`
	Content          nulls.String `col:"" json:"accountName" validate:"required" errm:"必填"`
	Snapshot         nulls.String `col:"" json:"snapshot"`
	CreateAt         nulls.String `col:"" json:"createAt"`
	Status           nulls.String `col:"" json:"status"`
	Memo             nulls.String `col:"" json:"memo"`
	Code             nulls.String `col:"" json:"code"`
	ApplicantUser_id nulls.Int    `col:"fk" json:"applicantUser_id"`
	ApproveUser_id   nulls.Int    `col:"fk" json:"approveUser_id"`

	//========fk
	ApplicantUser User `ref:"user,applicantUser_id" json:"applicantUser_id.row" validate:"-"`
	ApproveUser   User `ref:"user,approveUser_id" json:"approveUser_id.row" validate:"-"`
}

func (item *Application) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Content,
		&item.Snapshot,
		&item.CreateAt,
		&item.Status,
		&item.Memo,
		&item.Code,
		&item.ApplicantUser_id,
		&item.ApproveUser_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *Application) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkApplicant := User{}
	fkApprove := User{}

	columns = append(item.Receivers(), fkApplicant.Receivers()...)
	columns = append(columns, fkApprove.Receivers()...)

	err := r.Scan(columns...)

	item.ApplicantUser = fkApplicant
	item.ApproveUser = fkApprove

	return err
}

func (item *Application) ScanRows(r *sql.Rows) error {

	var columns []interface{}

	fkApplicant := User{}
	fkApprove := User{}

	columns = append(item.Receivers(), fkApplicant.Receivers()...)
	columns = append(columns, fkApprove.Receivers()...)

	err := r.Scan(columns...)

	item.ApplicantUser = fkApplicant
	item.ApproveUser = fkApprove

	return err
}
