package commandList

type AddCommand struct {
	UserId string
	Amount int
  }
  
  type QuoteCommand struct {
	UserId string
	StockSymbol string
   }
  
  type BuyCommand struct {
	UserId string
	StockSymbol string
	Amount int
  }
  
  type CommitBuyCommand struct {
	UserId string
  }
  
  type CancelBuyCommand struct {
	UserId string
  }
  
  type SellCommand struct {
	UserId string
	StockSymbol string
	Amount int
  }
  
  type CommitSellCommand struct {
	UserId string
  }
  
  type CancelSellCommand struct {
	UserId string
  }
  
  type SetBuyAmountCommand struct {
	UserId string
	StockSymbol string
	Amount int
  }
  
  type CancelSetBuyCommand struct {
	UserId string
	StockSymbol string
  }
  
  type SetBuyTriggerCommand struct {
	UserId string
	StockSymbol string
	Amount int
  }
  
  type SetSellAmountCommand struct {
	UserId string
	StockSymbol string
	Amount int
  }
  
  type CancelSetSellCommand struct {
	UserId string
	StockSymbol string
  }
  
  type SetSellTriggerCommand struct {
	UserId string
	StockSymbol string
	Amount int
  }
  
  type DumplogCommand struct {
	UserId string
	StockSymbol string
  }
  
  type DisplaySummary struct {
	UserId string
  }