package models

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
	Auth string `json:"auth"`
}
