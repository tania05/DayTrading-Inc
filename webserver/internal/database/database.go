package database

import (
	"errors"
	"webserver/internal/money"
)

type Holding interface {
	pay() error
	receive() error
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
var nextTransactionId int64

func init() {
	userMap = make(map[string]money.Money)
	stockMap = make(map[string]map[string]int)
  nextTransactionId = 1
}

func NewTransactionId() int64 {
  v := nextTransactionId
  nextTransactionId++
  return v
}

func (hold StockHolding) pay() error {
	if stockMap[hold.UserId] == nil {
		stockMap[hold.UserId] = make(map[string]int)
	}

	beforeStockAmount := stockMap[hold.UserId][hold.StockSymbol]
	if beforeStockAmount < hold.Amount {
		return errors.New("User does not have the required stock")
	}

	stockMap[hold.UserId][hold.StockSymbol] -= hold.Amount
	return nil
}

func (hold StockHolding) receive() error {
	if stockMap[hold.UserId] == nil {
		stockMap[hold.UserId] = make(map[string]int)
	}

	stockMap[hold.UserId][hold.StockSymbol] += hold.Amount
	return nil
}

func (hold MoneyHolding) pay() error {
	beforeMoneyAmount := userMap[hold.UserId]
	if beforeMoneyAmount < hold.Amount {
		return errors.New("User does not have the required money")
	}

	userMap[hold.UserId] -= hold.Amount
	return nil
}

func (hold MoneyHolding) receive() error {
	userMap[hold.UserId] += hold.Amount
	return nil
}

func AddFunds(userId string, amount money.Money) error {
	receivable := MoneyHolding{UserId: userId, Amount: amount}
	return receivable.receive()
}

func CheckFunds(userId string) (money.Money, error) {
	return userMap[userId], nil
}

func CheckStock(userId string, stockSymbol string) (int, error) {
	if stockMap[userId] == nil {
		stockMap[userId] = make(map[string]int)
	}
	return stockMap[userId][stockSymbol], nil
}

func attemptAllocate(trans Transaction) (Transaction, error) {
	err := trans.payable.pay()
	if err != nil {
		return Transaction{}, err
	}

	return trans, nil
}

func AllocateFunds(userId string, amount money.Money, stockSymbol string, stockAmount int) (Transaction, error) {
	payable := MoneyHolding{UserId: userId, Amount: amount}
	receivable := StockHolding{UserId: userId, StockSymbol: stockSymbol, Amount: stockAmount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(trans)
}

func AllocateStocks(userId string, stockSymbol string, stockAmount int, amount money.Money) (Transaction, error) {
	payable := StockHolding{UserId: userId, StockSymbol: stockSymbol, Amount: stockAmount}
	receivable := MoneyHolding{UserId: userId, Amount: amount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(trans)
}

func HoldStocks(userId string, stockSymbol string, amount int) (Holding, error) {
  hold := StockHolding{UserId: userId, StockSymbol: stockSymbol, Amount: amount}

  err := hold.pay()
  if err != nil {
    return nil, err
  }

  return hold, nil
}

func HoldMoney(userId string, amount money.Money) (Holding, error) {
  hold := MoneyHolding{UserId: userId, Amount: amount}

  err := hold.pay()
  if err != nil {
    return nil, err
  }

  return hold, nil
}

//TODO return and execute should be atomic
func Return(holds... Holding) error {
  executed := make([]Holding, len(holds))

  for _, hold := range holds {
    h := Holding(hold)
    err := h.receive()
    if err != nil {
      // When we add the database, we should just cancel
      for _, usedHold := range executed {
        usedHold.pay()
      }
      return err
    }
    executed = append(executed, h)
  }

  return nil
}

func Execute(holds... Holding) error {
  executed := make([]Holding, len(holds))

  for _, hold := range holds {
    h := Holding(hold)
    err := h.pay()
    if err != nil {
      // When we add the database, we should just cancel
      for _, usedHold := range executed {
        usedHold.receive()
      }
      return err
    }
    executed = append(executed, h)
  }

  return nil
}

func Commit(trans Transaction) error {
	return trans.receivable.receive()
}

func Cancel(trans Transaction) error {
	return trans.payable.receive()
}