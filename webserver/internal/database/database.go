package database

import (
  "webserver/internal/money"
)

type Holding int

func AddFunds(userId string, amount money.Money) error {
  return nil
}

func CheckFunds(userId string) (money.Money, error) {
  return 230, nil
}

func CheckStock(userId string, stockSymbol string) (int, error) {
  return 14, nil
}

func AllocateFunds(userId string, amount money.Money) (Holding, error) {
  return 1, nil
}

func AllocateStocks(userId string, stockSymbol string, amount int) (Holding, error) {
  return 1, nil
}

func Commit(holding Holding) error {
  return nil
}

func Cancel(holding Holding) error {
  return nil
}
