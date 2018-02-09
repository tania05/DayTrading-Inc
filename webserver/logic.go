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
	"webserver/internal/config"
)

var buystack map[string][]database.Transaction = make(map[string][]database.Transaction)
var sellstack map[string][]database.Transaction = make(map[string][]database.Transaction)

func getQuotetest(ctx *context.Context) money.Money{
	return money.Money(23)
}

func getQuote(ctx *context.Context) money.Money {

	addr := config.GlobalConfig.WebServer.Domain+ ":" + strconv.Itoa(config.GlobalConfig.WebServer.Port)

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
	//fmt.Println(ctx.UserId)
	//fmt.Println(ctx.Funds)
	return database.AddFunds(ctx, amount)
}

func transact(ctx *context.Context, bs int, amount money.Money) {
	price := getQuote(ctx)
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))

	if bs == 1 {
		t, err := database.AllocateFunds(ctx, cost, stocknum)
		// if err != nil {
		// 	panic(err)
		// }
		if err == nil {
			buystack[ctx.UserId] = append(buystack[ctx.UserId], t)
			fmt.Println("buy")
		}
		//err := pushPendingBuy(cost, stocknum, stock, user)
	} else {
		t, err := database.AllocateStocks(ctx, stocknum, cost)
		// if err != nil {
		// 	panic(err)
		// }
		if err == nil {
			sellstack[ctx.UserId] = append(sellstack[ctx.UserId], t)
			// fmt.Println(sellstack[ctx.UserId])
			// fmt.Println(t)
			// fmt.Println("sell")
		}
	}
}

func commitTransact(ctx *context.Context, bs int){
	//fmt.Println("buy/sell confirm")

	if bs == 1 && len(buystack[ctx.UserId]) > 0{
		database.Commit(ctx, popback(buystack[ctx.UserId]))
		fmt.Println("buy confirm")
	} else if len(sellstack[ctx.UserId]) > 0{
		database.Commit(ctx, popback(sellstack[ctx.UserId]))
		fmt.Println("sell confirm")
	}
}

func cancelTransact(ctx *context.Context, bs int){
	if bs == 1 && len(buystack[ctx.UserId]) > 0 {
		database.Cancel(ctx, popback(buystack[ctx.UserId]))
	} else if len(sellstack[ctx.UserId]) > 0{
		fmt.Println(sellstack[ctx.UserId])
		database.Cancel(ctx, popback(sellstack[ctx.UserId]))
	}
}

func popback(s []database.Transaction) database.Transaction {
	x, a := s[len(s)-1], s[:len(s)-1]
	s = a
	return x
}
