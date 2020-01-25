package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Company struct {
	ID                    int          `col:"" json:"id"`
	CompanyType           nulls.Int    `col:"" json:"companyType" validate:"required" errm:"公司类型不可为空"`
	Code                  nulls.String `col:"" json:"code"` // 权限可以后期改
	Name                  nulls.String `col:"" json:"name" validate:"required" errm:"公司名必填"`
	EName                 nulls.String `col:"" json:"ename"`
	ShortName             nulls.String `col:"" json:"shortname"`
	EShortName            nulls.String `col:"" json:"eshortname"`
	Address               nulls.String `col:"" json:"address"`
	Postcode              nulls.String `col:"" json:"postcode"`
	Phone1                nulls.String `col:"" json:"phone1"`
	Phone2                nulls.String `col:"" json:"phone2"`
	Phone3                nulls.String `col:"" json:"phone3"`
	Fax1                  nulls.String `col:"" json:"fax1"`
	Fax2                  nulls.String `col:"" json:"fax2"`
	Email1                nulls.String `col:"" json:"email1"`
	Email2                nulls.String `col:"" json:"email2"`
	Website               nulls.String `col:"" json:"website"`
	Memo                  nulls.String `col:"" json:"memo"`
	IsActive              nulls.Bool   `col:"" json:"isActive"`
	RetrieveTime          nulls.Time   `col:"" json:"retrieveTime"`
	UpdateAt              nulls.Time   `col:"" json:"updateAt"`
	CreateAt              nulls.Time   `col:"" json:"createAt"`
	Gsfj                  nulls.String `col:"" json:"gsfj"`
	Fjdz                  nulls.String `col:"" json:"fjdz"`
	Fjyb                  nulls.String `col:"" json:"fjyb"`
	TaxCode               nulls.String `col:"" json:"taxcode"`
	IsDelete              nulls.Bool   `col:"" json:"isDelete"`
	Retriever_id          nulls.Int    `col:"" json:"retriever_id"`
	UpdateUser_id         nulls.Int    `col:"" json:"updateUser_id"`
	Region_id             nulls.Int    `col:"" json:"region_id"`
	ImageLicense_id       nulls.Int    `col:"" json:"imageLicense_id"`
	ImageLicense          Image        `json:"imageLicense.row"`
	ImageBizCard_id       nulls.Int    `col:"" json:"imageBizCard_id"`
	ImageBizCard          Image        `json:"imageBizCard_id.row"`
	Gallary_folder_id     nulls.Int    `col:"" json:"gallary_folder_id"`
	Retriever_id_userName nulls.String `json:"retriever_id.userName"`
}

type CompanyList struct {
	Items []*Company
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry
func (item *Company) ScanRow(r *sql.Row) error {
	return r.Scan(
		&item.ID,
		&item.CompanyType,
		&item.Code,
		&item.Name,
		&item.EName,
		&item.ShortName,
		&item.EShortName,
		&item.Address,
		&item.Postcode,
		&item.Phone1,
		&item.Phone2,
		&item.Phone3,
		&item.Fax1,
		&item.Fax2,
		&item.Email1,
		&item.Email2,
		&item.Website,
		&item.Memo,
		&item.IsActive,
		&item.RetrieveTime,
		&item.UpdateAt,
		&item.CreateAt,
		&item.Gsfj,
		&item.Fjdz,
		&item.Fjyb,
		&item.TaxCode,
		&item.IsDelete,
		&item.Retriever_id,
		&item.UpdateUser_id,
		&item.Region_id,
		&item.ImageLicense_id,
		&item.ImageBizCard_id,
		&item.Gallary_folder_id)
}

func (item *Company) ScanRows(r *sql.Rows) error {
	return r.Scan(
		&item.ID,
		&item.CompanyType,
		&item.Code,
		&item.Name,
		&item.EName,
		&item.ShortName,
		&item.EShortName,
		&item.Address,
		&item.Postcode,
		&item.Phone1,
		&item.Phone2,
		&item.Phone3,
		&item.Fax1,
		&item.Fax2,
		&item.Email1,
		&item.Email2,
		&item.Website,
		&item.Memo,
		&item.IsActive,
		&item.RetrieveTime,
		&item.UpdateAt,
		&item.CreateAt,
		&item.Gsfj,
		&item.Fjdz,
		&item.Fjyb,
		&item.TaxCode,
		&item.IsDelete,
		&item.Retriever_id,
		&item.UpdateUser_id,
		&item.Region_id,
		&item.ImageLicense_id,
		&item.ImageBizCard_id,
		&item.Gallary_folder_id,
		&item.Retriever_id_userName)
}

func (list *CompanyList) ScanRow(r *sql.Rows) error {

	item := new(Company) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
