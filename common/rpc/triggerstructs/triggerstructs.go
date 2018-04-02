package triggerstructs

const FSetBuyAmountCommand = "SetBuyAmountCommand"
type SetBuyAmountCommand struct {
	TransactionNum int64
	UserId string
	StockSymbol string
	Amount int
}

const FSetBuyTriggerCommand = "SetBuyTriggerCommand"
type SetBuyTriggerCommand struct {
	TransactionNum int64
	UserId string
	StockSymbol string
	ExecutionPrice int
}

const FSetSellAmountCommand = "SetBuyAmountCommand"
type SetSellAmountCommand struct {
	TransactionNum int64
	UserId string
	StockSymbol string
	Amount int

}

const FSetSellTriggerCommand = "SetSellTriggerCommand"
type SetSellTriggerCommand struct {
	TransactionNum int64
	UserId string
	StockSymbol string
	ExecutionPrice int
}

const FCancelSetBuyCommand = "CancelSetBuyCommand"
type CancelSetBuyCommand struct {
	TransactionNum int64
	UserId string
	StockSymbol string
}

const FCancelSetSellCommand = "CancelSetSellCommand"
type CancelSetSellCommand struct {
	TransactionNum int64
	UserId string
	StockSymbol string
}
