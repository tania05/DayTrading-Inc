package database

import (
	"webserver/internal/money"
	"webserver/internal/logger"
	"webserver/internal/context"
	"sync"

	//"database/sql"
	//_ "github.com/lib/pq"
	//"fmt"
)

type Holding interface {
	pay(*context.Context) error
	receive(*context.Context) error
}

type Transaction struct {
	payable    Holding
	receivable Holding
}

type StockHolding struct {
	UserId      string
	StockSymbol string
	Amount      int
}

type MoneyHolding struct {
	UserId string
	Amount money.Money
}

var userMap map[string]money.Money
var stockMap map[string]map[string]int
var mutex sync.Mutex

//func attemptOpenConnection() error {
//	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s"))
//}

func init() {
	userMap = make(map[string]money.Money)
	stockMap = make(map[string]map[string]int)
}

func (hold StockHolding) pay(ctx *context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	if stockMap[hold.UserId] == nil {
		stockMap[hold.UserId] = make(map[string]int)
	}

	beforeStockAmount := stockMap[hold.UserId][hold.StockSymbol]
	if beforeStockAmount < hold.Amount {
	  return ctx.MakeError("Not enough stocks to complete transaction")
	}

	stockMap[hold.UserId][hold.StockSymbol] -= hold.Amount
	return nil
}

func (hold StockHolding) receive(ctx *context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	if stockMap[hold.UserId] == nil {
		stockMap[hold.UserId] = make(map[string]int)
	}

	stockMap[hold.UserId][hold.StockSymbol] += hold.Amount
	return nil
}

func (hold MoneyHolding) pay(ctx *context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	beforeMoneyAmount := userMap[hold.UserId]
	ctx.Funds = beforeMoneyAmount

	if beforeMoneyAmount < hold.Amount {
	  return ctx.MakeError("Not enough money to complete transaction")
	}

	userMap[hold.UserId] -= hold.Amount
	ctx.Funds = userMap[hold.UserId]
	ctx.MakeAccountTransactionLog(logger.RemoveAction)

	return nil
}

func (hold MoneyHolding) receive(ctx *context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	userMap[hold.UserId] += hold.Amount
	ctx.Funds = userMap[hold.UserId]
	ctx.MakeAccountTransactionLog(logger.AddAction)
	return nil
}

func AddFunds(ctx *context.Context, amount money.Money) error {
  receivable := MoneyHolding{UserId: ctx.UserId, Amount: amount}
  err := receivable.receive(ctx)
  return err
}

func CheckFunds(ctx *context.Context) (money.Money, error) {
	ctx.Funds = userMap[ctx.UserId]
	return userMap[ctx.UserId], nil
}

func CheckStock(ctx *context.Context) (int, error) {
	if stockMap[ctx.UserId] == nil {
		stockMap[ctx.UserId] = make(map[string]int)
	}
	ctx.Funds = userMap[ctx.UserId]
	return stockMap[ctx.UserId][ctx.StockSymbol], nil
}

func attemptAllocate(ctx *context.Context, trans Transaction) (Transaction, error) {
	err := trans.payable.pay(ctx)
	if err != nil {
		return Transaction{}, err
	}

	return trans, nil
}

func AllocateFunds(ctx *context.Context, amount money.Money, stockAmount int) (Transaction, error) {
  payable := MoneyHolding{UserId: ctx.UserId, Amount: amount}
  receivable := StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: stockAmount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(ctx, trans)
}

func AllocateStocks(ctx *context.Context, stockAmount int, amount money.Money) (Transaction, error) {
  payable := StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: stockAmount}
  receivable := MoneyHolding{UserId: ctx.UserId, Amount: amount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(ctx, trans)
}

func HoldStocks(ctx *context.Context, amount int) (Holding, error) {
  hold := StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: amount}

  err := hold.pay(ctx)
  if err != nil {
    return nil, err
  }

  return hold, nil
}

func HoldMoney(ctx *context.Context, amount money.Money) (Holding, error) {
  hold := MoneyHolding{UserId: ctx.UserId, Amount: amount}

  err := hold.pay(ctx)
  if err != nil {
    return nil, err
  }

  return hold, nil
}

//TODO return and execute should be atomic
func Return(ctx *context.Context, holds... Holding) error {
  executed := make([]Holding, len(holds))

  for _, hold := range holds {
    h := Holding(hold)
    err := h.receive(ctx)
    if err != nil {
      // When we add the database, we should just cancel
      for _, usedHold := range executed {
        usedHold.pay(ctx)
      }
      return err
    }
    executed = append(executed, h)
  }

  return nil
}

func Execute(ctx *context.Context, holds... Holding) error {
  executed := make([]Holding, len(holds))

  for _, hold := range holds {
    h := Holding(hold)
    err := h.pay(ctx)
    if err != nil {
      // When we add the database, we should just cancel
      for _, usedHold := range executed {
        usedHold.receive(ctx)
      }
      return err
    }
    executed = append(executed, h)
  }

  return nil
}

func Commit(ctx *context.Context, trans Transaction) error {
	return trans.receivable.receive(ctx)
}

func Cancel(ctx *context.Context, trans Transaction) error {
	return trans.payable.receive(ctx)
}
