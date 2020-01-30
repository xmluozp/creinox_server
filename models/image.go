package models

import (
	"database/sql"

	"github.com/gobuffalo/nulls"
)

type Image struct {
	ID                  int          `col:"" json:"id"`
	Name                nulls.String `col:"" json:"name" validate:"required" errm:"角色名必填"`
	Height              nulls.Int    `col:"" json:"height"` // 权限可以后期改
	Width               nulls.Int    `col:"" json:"width" validate:"required" errm:"必须选择"`
	Sort                nulls.Int    `col:"" json:"sort"`
	Path                nulls.String `col:"" json:"path"`
	ThumbnailPath       nulls.String `col:"" json:"thumbnailPath"`
	Ext                 nulls.String `col:"" json:"ext"`
	CreateAt            nulls.String `col:"default" json:"createAt"`
	Gallary_folder_id   nulls.Int    `col:"" json:"gallary_folder_id"`
	Gallary_folder_memo nulls.Nulls  `json:"gallary_folder_id.memo"`
}

type ImageList struct {
	Items []*Image
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错

// learned from: https://stackoverflow.com/questions/53175792/how-to-make-scanning-db-rows-in-go-dry

func (item *Image) ScanRow(r *sql.Row) error {
	return r.Scan(
		&item.ID,
		&item.Name,
		&item.Height,
		&item.Width,
		&item.Sort,
		&item.Path,
		&item.ThumbnailPath,
		&item.Ext,
		&item.CreateAt,
		&item.Gallary_folder_id)
}

func (item *Image) ScanRows(r *sql.Rows) error {
	return r.Scan(
		&item.ID,
		&item.Name,
		&item.Height,
		&item.Width,
		&item.Sort,
		&item.Path,
		&item.ThumbnailPath,
		&item.Ext,
		&item.CreateAt,
		&item.Gallary_folder_id)
}

func (list *ImageList) ScanRow(r *sql.Rows) error {

	item := new(Image) // ---------- item

	if err := item.ScanRows(r); err != nil {
		return err
	}
	list.Items = append(list.Items, item)
	return nil
}
