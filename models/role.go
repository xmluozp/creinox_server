package models

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required" errm:"角色名必填"`
	Rank int    `json:"rank"`
	Auth string `json:"auth" validate:"required" errm:"必须选择"`
}

// 取的时候，类型[]byte就不关心是不是null。不然null转其他的报错
