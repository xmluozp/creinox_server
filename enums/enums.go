package enums

type contractType struct {
	SellContract  int
	BuyContract   int
	MouldContract int
}

var ContractType = contractType{
	SellContract:  1,
	BuyContract:   2,
	MouldContract: 3}

var ContractTypeLabel = map[int]string{
	1: "销售合同",
	2: "采购合同",
	3: "产品开发合同"}

//-------------------

type financialSubjectType struct {
	Payable    string
	Receivable string
}

var FinancialSubjectType = financialSubjectType{
	Payable:    "应付款",
	Receivable: "应收款"}
