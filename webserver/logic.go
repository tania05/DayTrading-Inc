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
	"github.com/go-redis/redis"
	"github.com/go-redis/cache"
	"encoding/json"
	"sync"
)

type QuoteCacheItem struct {
	Price money.Money
	Timestamp int64
}

var buystack map[string][]database.Transaction = make(map[string][]database.Transaction)
var sellstack map[string][]database.Transaction = make(map[string][]database.Transaction)
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
			return nil
		}

		redisCodec = &cache.Codec{
			Redis: redisClient,
			Marshal: json.Marshal,
			Unmarshal: json.Unmarshal,
		}
	}

	return redisCodec
}

func getQuotetest(ctx *context.Context) money.Money{
	return money.Money(23)
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
	return database.AddFunds(ctx, amount)
}

func transact(ctx *context.Context, bs int, amount money.Money) {
	price := getQuote(ctx)
	stocknum := int((amount/price))
	cost := money.Money(stocknum * int(price))
	
	mutex.Lock()
	defer mutex.Unlock()

	if bs == 1 {
		t, err := database.AllocateFunds(ctx, cost, stocknum)
		// if err != nil {
		// 	panic(err)
		// }
		if err == nil {
			
			buystack[ctx.UserId] = append(buystack[ctx.UserId], t)
			// fmt.Println("buy")
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
	mutex.Lock()
	defer mutex.Unlock()
	if bs == 1 && len(buystack[ctx.UserId]) > 0{
		database.Commit(ctx, popback(buystack[ctx.UserId]))
		// fmt.Println("buy confirm")
	} else if len(sellstack[ctx.UserId]) > 0{
		database.Commit(ctx, popback(sellstack[ctx.UserId]))
		// fmt.Println("sell confirm")
	}
}

func cancelTransact(ctx *context.Context, bs int){
	mutex.Lock()
	defer mutex.Unlock()
	if bs == 1 && len(buystack[ctx.UserId]) > 0 {
		database.Cancel(ctx, popback(buystack[ctx.UserId]))
	} else if len(sellstack[ctx.UserId]) > 0{
		// fmt.Println(sellstack[ctx.UserId])
		database.Cancel(ctx, popback(sellstack[ctx.UserId]))
	}
}

func popback(s []database.Transaction) database.Transaction {
	x, a := s[len(s)-1], s[:len(s)-1]
	s = a
	return x
}
