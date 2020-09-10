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

var LogActions = map[string]string{
	"c": "创建",
	"u": "修改",
	"d": "删除"}

//-------------------

type financialLedgerType struct {
	UnDecided     int
	PayableDebit  int
	PayableCredit int

	ReceivableDebit  int
	ReceivableCredit int

	PayablePayDebit  int
	PayablePayCredit int

	ReceivablePayDebit  int
	ReceivablePayCredit int
}

var FinancialLedgerType = financialLedgerType{
	UnDecided:           6,
	PayableDebit:        11,
	PayableCredit:       8,
	ReceivableDebit:     7,
	ReceivableCredit:    9,
	PayablePayDebit:     8,
	PayablePayCredit:    10,
	ReceivablePayDebit:  10,
	ReceivablePayCredit: 7}
