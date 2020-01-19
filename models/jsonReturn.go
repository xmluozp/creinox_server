package models

type JsonRowsReturn struct {
	Pagination  interface{}       `json:"pagination"`
	SearchTerms map[string]string `json:"searchTerms"`
	Rows        interface{}       `json:"rows"`
	Row         interface{}       `json:"row"`
	Message     map[string]string `json:"message"` // 所有message都是string。这个是显示在表单的
	Info        string            `json:"info"`    // 这个是显示在弹出框的（如果有的话）
}
