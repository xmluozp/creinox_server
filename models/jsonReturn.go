package models

type JsonRowsReturn struct {
	Pagination  interface{}       `json:"pagination"`
	SearchTerms map[string]string `json:"searchTerms"`
	Rows        interface{}       `json:"rows"`
	Row         interface{}       `json:"row"`
	Message     map[string]string `json:"message"` // 所有message都是string
}
