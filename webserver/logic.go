package main

import (
	"net"
	"log"
	"fmt"
	"strconv"
	"webserver/internal/database"
	"webserver/internal/money"
	"strings"
)

const domain = "localhost"
const port = 8080

var onestack database.Transaction

func getQuote(user string, stock string) money.Money {
	addr := (domain + ":" + strconv.Itoa(port))

	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	conn.Write([]byte(stock +","+user))
	conn.Write([]byte("\n"))

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Recieve: %s", buff[:n])
	val, err := strconv.Atoi(strings.Split(string(buff),",")[0])
	return money.Money(val)
}

func addFunds(user string, amount money.Money){
	database.AddFunds(user, amount)
	fmt.Println("add work?")
}

func transact(bs int, user string, amount money.Money, stock string) {
	price := money.Money(25)//getQuote(user, stock)
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))

	if bs == 1 {
		t, _ := database.AllocateFunds(user, cost, stock, stocknum)
		onestack = t
		fmt.Println("buy work?")
		//err := pushPendingBuy(cost, stocknum, stock, user)
	} else {
		t, _ := database.AllocateStocks(user, stock, stocknum, cost)
		onestack = t
		fmt.Println("sell work?")
	}
}

func commitTransact(bs int, user string){
	if bs == 1 {
		database.Commit(onestack)
	} else {
		database.Commit(onestack)
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