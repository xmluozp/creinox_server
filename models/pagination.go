package models

// 翻页用的
type Pagination struct {
	Page       int    `json:"page"`
	RowCount   int    `json:"rowCount"`
	PerPage    int    `json:"perPage"`
	TotalCount int    `json:"totalCount"`
	TotalPage  int    `json:"totalPage"`
	Order      string `json:"order"`
	OrderBy    string `json:"orderBy"`
}
