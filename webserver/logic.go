package main

import (
	"fmt"
	"time"
	"webserver/internal/transaction"
	"common/money"
	"common/context"
	"sync"
	"common/quote"
)

var mutex sync.Mutex

func addFunds(ctx *context.Context, amount money.Money ) error {
	//fmt.Println(ctx.UserId)
	//fmt.Println(ctx.Funds)
	return transaction.AddFunds(ctx, amount)
}

func transact(ctx *context.Context, bs int, amount money.Money) error {
	price, err := quote.GetQuote(ctx)
	if err != nil {
		return err
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
		return err
	}
	go func(ctx *context.Context, id int) {
		time.Sleep(60 * time.Second)
		fmt.Println("Checking if transaction ", id, " still exists, if so cancelling")
		transaction.CancelByTimeout(ctx, id)
	}(ctx, tx.Id)

	return nil
}

func commitTransact(ctx *context.Context, bs int){
	transaction.CommitTransaction(ctx, bs == 1)
}

func cancelTransact(ctx *context.Context, bs int){
	transaction.CancelTransaction(ctx, bs == 1)
}
