package main

import (
	"net"
	"log"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"webserver/internal/database"
	"webserver/internal/money"
	"webserver/internal/logger"
	"webserver/internal/context"
  "webserver/internal/trigger"
	"strings"
)

const domain = "localhost"
const port = 4441

var buystack []database.Transaction
var sellstack []database.Transaction

func getQuote(ctx *context.Context) money.Money {
	addr := (domain + ":" + strconv.Itoa(port))

	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
    	panic(err)
	}

	conn.Write([]byte(ctx.StockSymbol +","+ ctx.UserId))
	conn.Write([]byte("\n"))

	buff, _ := ioutil.ReadAll(conn)
	log.Printf("Recieve: %s", buff)

	f := strings.Split(string(buff), ",")
	f3, _:= strconv.Atoi(f[3])
	val, err := strconv.Atoi(strings.Replace(strings.Split(string(buff),",")[0],".","",-1))

	logger.Log(logger.QuoteServerLog{
		Timestamp: time.Now().UnixNano() / 1e6,
		Server: "ts0",
		TransactionNum: ctx.TransactionNum,
		QuoteServerTime: int64(f3),
		Username: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
		Price: money.Money(val),
		Cryptokey: strings.Trim(f[4], "\n")})
	fmt.Println(val)

	//TODO fix this horrible thingers //"kk", []sg if not empty
	trigger.OnQuoteUpdate(ctx, money.Money(val))
	return money.Money(val)
}

func addFunds(ctx *context.Context, amount money.Money ) error {
	return database.AddFunds(ctx, amount)
}

func transact(ctx *context.Context, bs int, amount money.Money) {
	price := getQuote(ctx)
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))

	if bs == 1 {
		t, err := database.AllocateFunds(ctx, cost, stocknum)
		if err != nil {
			panic(err)
		}
		buystack = append(buystack, t)
		fmt.Println("buy")
		//err := pushPendingBuy(cost, stocknum, stock, user)
	} else {
		t, _ := database.AllocateStocks(ctx, stocknum, cost)
		sellstack = append(sellstack, t)
		fmt.Println("sell")
	}
}

func commitTransact(ctx *context.Context, bs int){
	//fmt.Println("buy/sell confirm")

	if bs == 1 && len(buystack) > 0{
		database.Commit(ctx, popback(buystack))
		fmt.Println("buy confirm")
	} else if len(sellstack) > 0{
		database.Commit(ctx, popback(sellstack))
		fmt.Println("sell confirm")
	}
}

func cancelTransact(ctx *context.Context, bs int){
	if bs == 1 && len(buystack) > 0 {
		database.Cancel(ctx, popback(buystack))
	} else if len(sellstack) > 0{
		database.Cancel(ctx, popback(sellstack))
	}
}

func popback(s []database.Transaction) database.Transaction {
	x, a := s[len(s)-1], s[:len(s)-1]
	s = a
	return x
}
