package models

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/initial"
)

type Image struct {
	ID                nulls.Int    `col:"" json:"id"`
	Name              nulls.String `col:"" json:"name" validate:"required" errm:"必填"`
	Height            nulls.Int    `col:"" json:"height"` // 权限可以后期改
	Width             nulls.Int    `col:"" json:"width"`
	Sort              nulls.Int    `col:"" json:"sort"`
	Path              nulls.String `col:"" json:"path"`
	ThumbnailPath     nulls.String `col:"" json:"thumbnailPath"`
	Ext               nulls.String `col:"" json:"ext"`
	CreateAt          nulls.String `col:"default" json:"createAt"`
	Gallary_folder_id nulls.Int    `col:"" json:"gallary_folder_id"`

	Gallary_folder_memo nulls.String `json:"gallary_folder_id.memo"`

	FileName          nulls.String `col:"" json:"filename"`
	ThumbnailFileName nulls.String `col:"" json:"thumbnailfilename"`
}

type ImageList struct {
	Items []*Image
}

func (item *Image) Receivers() (itemPtrs []interface{}) {

	values := []interface{}{
		&item.ID,
		&item.Name,
		&item.Height,
		&item.Width,
		&item.Sort,
		&item.Path,
		&item.ThumbnailPath,
		&item.Ext,
		&item.CreateAt,
		&item.Gallary_folder_id}

	valuePtrs := make([]interface{}, len(values))

	for i := range values {
		valuePtrs[i] = values[i]
	}

	return valuePtrs
}

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *Image) ScanRow(r *sql.Row) error {

	err := r.Scan(append(item.Receivers(), &item.Gallary_folder_memo)...)

	return err
}

func (item *Image) ScanRows(r *sql.Rows) error {

	err := r.Scan(append(item.Receivers(), &item.Gallary_folder_memo)...)

	return err
}

func (item *Image) Getter() Image {

	// cfg, err := goconfig.LoadConfigFile("conf.ini")

	// if err != nil {
	// 	panic("错误，找不到conf.ini配置文件")
	// }

	// rootUrl, err := cfg.GetValue("site", "root")
	// port, err := cfg.Int("site", "port")
	// uploads, err := cfg.GetValue("site", "uploads")

	rootUrl, _, port, uploads := initial.GetConfig()

	uploadFolder := fmt.Sprintf("%s:%d/%s/", rootUrl, port, uploads)

	item.ThumbnailFileName = item.ThumbnailPath

	if item.ThumbnailPath.Valid {
		// item.ThumbnailPath = nulls.NewString(uploads + "/" + item.ThumbnailPath.String) // 20200518 因为cors从服务端解决所以去掉这个
		item.ThumbnailPath = nulls.NewString(uploadFolder + item.ThumbnailPath.String)
	}

	item.FileName = item.Path

	if item.Path.Valid {
		item.Path = nulls.NewString(uploadFolder + item.Path.String)
	}

	return *item
}

// 根据网站目录+文件名，生成一个长path，用来前端访问用
func (item *Image) AddPath(path string) string {

	if path == "" {
		return path
	}

	// cfg, err := goconfig.LoadConfigFile("conf.ini")

	// if err != nil {
	// 	panic("错误，找不到conf.ini配置文件")
	// }

	// rootUrl, err := cfg.GetValue("site", "root")
	// port, err := cfg.Int("site", "port")
	// uploads, err := cfg.GetValue("site", "uploads")

	rootUrl, _, port, uploads := initial.GetConfig()

	uploadFolder := fmt.Sprintf("%s:%d/%s/", rootUrl, port, uploads)

	return uploadFolder + path
}

func (list *ImageList) ScanRow(r *sql.Rows) error {

	item := new(Image) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
