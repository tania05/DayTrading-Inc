package main

import (
	"fmt"
	"time"
	"webserver/internal/transaction"
	"common/money"
	"common/context"
	"common/quote"
)

func addFunds(ctx *context.Context, amount money.Money ) (int, error) {
	//fmt.Println(ctx.UserId)
	//fmt.Println(ctx.Funds)
	return transaction.AddFunds(ctx, amount)
}

func transact(ctx *context.Context, bs int, amount money.Money) (money.Money, int, error) {
	price, err := quote.GetQuote(ctx)
	if err != nil {
		return -1, -1, err
	}
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))

	fmt.Println("Price ", price, " stocknum ", stocknum, " cost ", cost)

	var tx transaction.Transaction
	if bs == 1 {
		fmt.Println("Allocating funds")
		tx, err = transaction.AllocateFunds(ctx, cost, stocknum)
	} else {
		tx, err = transaction.AllocateStocks(ctx, stocknum, cost)
	}

	if err != nil {
		return -1, -1, err
	}
	go func(ctx *context.Context, id int) {
		time.Sleep(60 * time.Second)
		fmt.Println("Checking if transaction ", id, " still exists, if so cancelling")
		transaction.CancelByTimeout(ctx, id)
	}(ctx, tx.Id)

	return cost, stocknum, nil
}

func commitTransact(ctx *context.Context, bs int) error {
	return transaction.CommitTransaction(ctx, bs == 1)
}

func cancelTransact(ctx *context.Context, bs int) error {
	return transaction.CancelTransaction(ctx, bs == 1)
}
