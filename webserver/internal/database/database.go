package database

import (
	"errors"
	"webserver/internal/money"
)

type holding interface {
	pay() error
	receive() error
}

type Transaction struct {
	payable    holding
	receivable holding
}

type stockHolding struct {
	userId      string
	stockSymbol string
	amount      int
}

type moneyHolding struct {
	userId string
	amount money.Money
}

var userMap map[string]money.Money
var stockMap map[string]map[string]int

func init() {
	userMap = make(map[string]money.Money)
	stockMap = make(map[string]map[string]int)
}

func (hold stockHolding) pay() error {
	if stockMap[hold.userId] == nil {
		stockMap[hold.userId] = make(map[string]int)
	}

	beforeStockAmount := stockMap[hold.userId][hold.stockSymbol]
	if beforeStockAmount < hold.amount {
		return errors.New("User does not have the required stock")
	}

	stockMap[hold.userId][hold.stockSymbol] -= hold.amount
	return nil
}

func (hold stockHolding) receive() error {
	if stockMap[hold.userId] == nil {
		stockMap[hold.userId] = make(map[string]int)
	}

	stockMap[hold.userId][hold.stockSymbol] += hold.amount
	return nil
}

func (hold moneyHolding) pay() error {
	beforeMoneyAmount := userMap[hold.userId]
	if beforeMoneyAmount < hold.amount {
		return errors.New("User does not have the required money")
	}

	userMap[hold.userId] -= hold.amount
	return nil
}

func (hold moneyHolding) receive() error {
	userMap[hold.userId] += hold.amount
	return nil
}

func AddFunds(userId string, amount money.Money) error {
	receivable := moneyHolding{userId: userId, amount: amount}
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
	payable := moneyHolding{userId: userId, amount: amount}
	receivable := stockHolding{userId: userId, stockSymbol: stockSymbol, amount: stockAmount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(trans)
}

func AllocateStocks(userId string, stockSymbol string, stockAmount int, amount money.Money) (Transaction, error) {
	payable := stockHolding{userId: userId, stockSymbol: stockSymbol, amount: stockAmount}
	receivable := moneyHolding{userId: userId, amount: amount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(trans)
}

func Commit(trans Transaction) error {
	return trans.receivable.receive()
}

func Cancel(trans Transaction) error {
	return trans.payable.receive()
}
