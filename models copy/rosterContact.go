package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type RosterContact struct {
	ID         nulls.Int    `col:"" json:"id"`
	FullName   nulls.String `col:"" json:"fullName" validate:"required" errm:"必填"`
	EFullName  nulls.String `col:"" json:"eFullName"`
	Phone1     nulls.String `col:"" json:"phone1"`
	Phone2     nulls.String `col:"" json:"phone2"`
	Skype      nulls.String `col:"" json:"skype"`
	Email      nulls.String `col:"" json:"email"`
	Wechat     nulls.String `col:"" json:"wechat"`
	Whatsapp   nulls.String `col:"" json:"whatsapp"`
	Facebook   nulls.String `col:"" json:"facebook"`
	Memo       nulls.String `col:"" json:"memo"`
	Company_id nulls.Int    `col:"fk" json:"company_id"`
}

type RosterContactList struct {
	Items []*RosterContact
}

func (item *RosterContact) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.FullName,
		&item.EFullName,
		&item.Phone1,
		&item.Phone2,
		&item.Skype,
		&item.Email,
		&item.Wechat,
		&item.Whatsapp,
		&item.Facebook,
		&item.Memo,
		&item.Company_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

func (item *RosterContact) ScanRow(r *sql.Row) error {
	err := r.Scan(item.Receivers()...)

	return err
}

func (item *RosterContact) ScanRows(r *sql.Rows) error {
	err := r.Scan(item.Receivers()...)

	return err
}

func (list *RosterContactList) ScanRow(r *sql.Rows) error {

	item := new(RosterContact) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
