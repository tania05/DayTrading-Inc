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
	"strings"
)

const domain = "localhost"
const port = 4441

var onestack database.Transaction

func getQuote(user string, stock string) money.Money {
	addr := (domain + ":" + strconv.Itoa(port))

	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
    panic(err)
	}

	conn.Write([]byte(stock +","+user))
	conn.Write([]byte("\n"))

	buff, _ := ioutil.ReadAll(conn)
	log.Printf("Recieve: %s", buff)

	f := strings.Split(string(buff), ",")
	f3, _:= strconv.Atoi(f[3])
	val, err := strconv.Atoi(strings.Replace(strings.Split(string(buff),",")[0],".","",-1))

	logger.Log(logger.QuoteServerLog{
		Timestamp: time.Now().UnixNano() / 1e6,
		Server: "ts0",
		TransactionNum: 69,
		QuoteServerTime: int64(f3),
		Username: user,
		StockSymbol: stock,
		Price: money.Money(val),
		Cryptokey: f[4]})
	fmt.Println(val)
	return money.Money(val)
}

func addFunds(user string, amount money.Money) (int64, error) {
	return database.AddFunds(user, amount)
}

func transact(bs int, user string, amount money.Money, stock string) {
	price := money.Money(45)
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))

	if bs == 1 {
		t, err := database.AllocateFunds(user, cost, stock, stocknum)
		if err != nil {
			panic(err)
		}
		onestack = t
		fmt.Println("buy")
		//err := pushPendingBuy(cost, stocknum, stock, user)
	} else {
		t, _ := database.AllocateStocks(user, stock, stocknum, cost)
		onestack = t
		fmt.Println("sell")
	}
}

func commitTransact(bs int, user string){
	//fmt.Println("buy/sell confirm")

	if bs == 1 {
		database.Commit(onestack)
		fmt.Println("buy confirm")

	} else {
		database.Commit(onestack)
		fmt.Println("sell confirm")

	}

}

func cancelTransact(bs int, user string){
	if bs == 1 {
		database.Cancel(onestack)
	} else {
		database.Cancel(onestack)
	}
}

// func setTransAmount(bs int, user string, stock string, amount money){
// 	if bs == 1 {
// 		err = addBuyAmount(user, stock, amount)
// 	}
// 	else {
// 		err = addSellAmount(user, stock, amount)
// 	}
// }

// func setTransTrigger(bs int, user string, stock string, amount money){
// 	if bs == 1 {
// 		err = setBuyTrigger(user, stock, amount)
// 	}
// 	else {
// 		err = setSellTrigger(user, stock, amount)
// 	}
// }

// func cancelTransSet(bs int, user string, stock string){
// 	if bs == 1 {
// 		err = cancelBuyTrigger(user, stock)
// 	}
// 	else {
// 		err = cancelSellTrigger(user, stock)
// 	}
// }


// func main() {

// 	c := getQuote("TST,user")
// 	log.Printf("%s", c)
// }
