package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Company struct {
	ID           nulls.Int    `col:"" json:"id"`
	CompanyType  nulls.Int    `col:"" json:"companyType" validate:"required" errm:"公司类型不可为空"`
	Code         nulls.String `col:"" json:"code"` // 权限可以后期改
	Name         nulls.String `col:"" json:"name" validate:"required" errm:"公司名必填"`
	EName        nulls.String `col:"" json:"ename"`
	ShortName    nulls.String `col:"" json:"shortname"`
	EShortName   nulls.String `col:"" json:"eshortname"`
	Address      nulls.String `col:"" json:"address"`
	EAddress     nulls.String `col:"" json:"eaddress"`
	Postcode     nulls.String `col:"" json:"postcode"`
	Phone1       nulls.String `col:"" json:"phone1"`
	Phone2       nulls.String `col:"" json:"phone2"`
	Phone3       nulls.String `col:"" json:"phone3"`
	Fax1         nulls.String `col:"" json:"fax1"`
	Fax2         nulls.String `col:"" json:"fax2"`
	Email1       nulls.String `col:"" json:"email1"`
	Email2       nulls.String `col:"" json:"email2"`
	Website      nulls.String `col:"" json:"website"`
	Memo         nulls.String `col:"" json:"memo"`
	IsActive     nulls.Bool   `col:"" json:"isActive"`
	RetrieveTime nulls.Time   `col:"" json:"retrieveTime"`
	UpdateAt     nulls.Time   `col:"newtime" json:"updateAt"`
	CreateAt     nulls.Time   `col:"default" json:"createAt"`
	Gsfj         nulls.String `col:"" json:"gsfj"`
	Fjdz         nulls.String `col:"" json:"fjdz"`
	Fjyb         nulls.String `col:"" json:"fjyb"`
	TaxCode      nulls.String `col:"" json:"taxcode"`
	IsDelete     nulls.Bool   `col:"" json:"isDelete"`

	// 内部公司专用字段-----
	Zsl nulls.Float32 `col:"" json:"zsl"` //增税率
	Hl  nulls.Float32 `col:"" json:"hl"`  //汇率
	Tsl nulls.Float32 `col:"" json:"tsl"` // 退税率
	// 内部公司专用字段-----

	// 搜索用
	KeyWord nulls.String `json:"keyword" keywords:"code|name|ename|shortname|eshortname"`

	Retriever_id          nulls.Int    `col:"fk" json:"retriever_id"`
	UpdateUser_id         nulls.Int    `col:"fk" json:"updateUser_id"`
	Gallary_folder_id     nulls.Int    `col:"" json:"gallary_folder_id"` // no fk constraint here
	ImageLicense_id       nulls.Int    `col:"fk" json:"imageLicense_id,omitempty"`
	ImageBizCard_id       nulls.Int    `col:"fk" json:"imageBizCard_id,omitempty"`
	Region_id             nulls.Int    `col:"" json:"region_id"`
	Retriever_id_userName nulls.String `json:"retriever_id.userName"`

	ImageLicense Image `ref:"image,imageLicense_id" json:"imageLicense_id.row" validate:"-"`
	ImageBizCard Image `ref:"image,imageBizCard_id" json:"imageBizCard_id.row" validate:"-"`

	RetrieverItem User `ref:"user,retriever_id" json:"retriever_id.row" validate:"-"`
}

type CompanyList struct {
	Items []*Company
}

func (item *Company) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
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
		&item.Gallary_folder_id,
		&item.ImageLicense_id,
		&item.ImageBizCard_id,
		&item.Region_id,
		&item.EAddress,
		&item.Zsl,
		&item.Hl,
		&item.Tsl}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry
func (item *Company) ScanRow(r *sql.Row) error {

	var columns []interface{}

	fkImageLicense := Image{}
	fkImageBizCard := Image{}
	fkRetriever := User{}

	columns = append(item.Receivers(), fkImageLicense.Receivers()...)
	columns = append(columns, fkImageBizCard.Receivers()...)
	columns = append(columns, fkRetriever.Receivers()...)

	err := r.Scan(columns...)

	item.ImageLicense = fkImageLicense.Getter()
	item.ImageBizCard = fkImageBizCard.Getter()

	item.RetrieverItem = fkRetriever

	return err
}

func (item *Company) ScanRows(r *sql.Rows) error {
	// err := r.Scan(item.Receivers()...)

	// return err

	var columns []interface{}

	fkImageLicense := Image{}
	fkImageBizCard := Image{}
	fkRetriever := User{}

	columns = append(item.Receivers(), fkImageLicense.Receivers()...)
	columns = append(columns, fkImageBizCard.Receivers()...)
	columns = append(columns, fkRetriever.Receivers()...)

	err := r.Scan(columns...)

	item.RetrieverItem = fkRetriever

	return err
}

func (list *CompanyList) ScanRow(r *sql.Rows) error {

	item := new(Company) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
