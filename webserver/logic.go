package main

import (
	"net"
	"log"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"webserver/internal/transaction"
	"webserver/internal/money"
	"webserver/internal/logger"
	"webserver/internal/context"
  "webserver/internal/trigger"
	"strings"
	"common/config"
	"github.com/go-redis/redis"
	"github.com/go-redis/cache"
	"encoding/json"
	"sync"
)

type QuoteCacheItem struct {
	Price money.Money
	Timestamp int64
}

var mutex sync.Mutex

var redisCodec *cache.Codec

func getClient() *cache.Codec {
	addr := config.GlobalConfig.Redis.Domain + ":" + strconv.Itoa(config.GlobalConfig.Redis.Port)

	if redisCodec == nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr: addr,
			Password: "",
			DB: 0,
		})

		pingError := redisClient.Ping().Err()
		if pingError != nil {
			fmt.Println("Failed to contact Redis server")
			fmt.Println(pingError)
			panic(pingError)
			return nil
		}
		fmt.Println("Connected to Redis")

		redisCodec = &cache.Codec{
			Redis: redisClient,
			Marshal: json.Marshal,
			Unmarshal: json.Unmarshal,
		}
	}

	return redisCodec
}

func getQuote(ctx *context.Context) money.Money {
	codec := getClient()

	var cacheVal QuoteCacheItem
	if codec != nil {
		if err := codec.Get(ctx.StockSymbol, &cacheVal); err == nil {
			// todo we already have this value stored, systemeventlog or something
			fmt.Println("Got price from cache")
			fmt.Println(cacheVal)
			return cacheVal.Price
		}
	}

	addr := config.GlobalConfig.QuoteServer.Domain+ ":" + strconv.Itoa(config.GlobalConfig.QuoteServer.Port)

	fmt.Println("Contacting " + addr)
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

	now := time.Now().UnixNano() / 1e6

	logger.Log(logger.QuoteServerLog{
		Timestamp: now,
		Server: "ts0",
		TransactionNum: ctx.TransactionNum,
		QuoteServerTime: int64(f3),
		Username: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
		Price: money.Money(val),
		Cryptokey: strings.Trim(f[4], "\n")})
	fmt.Println(val)


	cacheItem := QuoteCacheItem {
		Price: money.Money(val),
		Timestamp: now,
	}

	err = codec.Set(&cache.Item{
		Key: ctx.StockSymbol,
		Object: cacheItem,
		Expiration: time.Second * 45,
	})
	if err != nil {
		fmt.Println("Failed to set cache key")
		fmt.Println(err)
	}

	trigger.OnQuoteUpdate(ctx, money.Money(val))
	return money.Money(val)
}

func addFunds(ctx *context.Context, amount money.Money ) error {
	//fmt.Println(ctx.UserId)
	//fmt.Println(ctx.Funds)
	return transaction.AddFunds(ctx, amount)
}

func transact(ctx *context.Context, bs int, amount money.Money) error {
	price := getQuote(ctx)
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))

	fmt.Println("Price ", price, " stocknum ", stocknum, " cost ", cost)

	var tx transaction.Transaction
	var err error
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
