package quote

import (
	"common/money"
	"fmt"
	"common/config"
	"strconv"
	"net"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"common/logger"
	"github.com/go-redis/cache"
	"common/context"
	"github.com/go-redis/redis"
	"encoding/json"
)

type QuoteCacheItem struct {
	Price     money.Money
	Timestamp int64
}

var redisCodec *cache.Codec

func getClient() *cache.Codec {
	addr := config.GlobalConfig.Redis.Domain + ":" + strconv.Itoa(config.GlobalConfig.Redis.Port)

	if redisCodec == nil {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
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
			Redis:     redisClient,
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		}
	}

	return redisCodec
}

func GetQuote(ctx *context.Context) (money.Money, error) {
	codec := getClient()

	var cacheVal QuoteCacheItem
	if codec != nil {
		if err := codec.Get(ctx.StockSymbol, &cacheVal); err == nil {
			// todo we already have this value stored, systemeventlog or something
			fmt.Println("Got price from cache")
			fmt.Println(cacheVal)
			return cacheVal.Price, nil
		}
	}

	addr := config.GlobalConfig.QuoteServer.Domain + ":" + strconv.Itoa(config.GlobalConfig.QuoteServer.Port)

	fmt.Println("Contacting " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Quote server connection error ", err)
		return 0, err
	}

	defer conn.Close()


	conn.Write([]byte(ctx.StockSymbol + "," + ctx.UserId))
	conn.Write([]byte("\n"))

	buff, _ := ioutil.ReadAll(conn)
	log.Printf("Recieve: %s", buff)

	f := strings.Split(string(buff), ",")
	f3, _ := strconv.Atoi(f[3])
	val, err := strconv.Atoi(strings.Replace(strings.Split(string(buff), ",")[0], ".", "", -1))

	now := time.Now().UnixNano() / 1e6

	logger.Log(logger.QuoteServerLog{
		Timestamp:       now,
		Server:          "ts0",
		TransactionNum:  ctx.TransactionNum,
		QuoteServerTime: int64(f3),
		Username:        ctx.UserId,
		StockSymbol:     ctx.StockSymbol,
		Price:           money.Money(val),
		Cryptokey:       strings.Trim(f[4], "\n")})
	fmt.Println(val)

	cacheItem := QuoteCacheItem{
		Price:     money.Money(val),
		Timestamp: now,
	}

	err = codec.Set(&cache.Item{
		Key:        ctx.StockSymbol,
		Object:     cacheItem,
		Expiration: time.Second * 45,
	})
	if err != nil {
		fmt.Println("Failed to set cache key")
		fmt.Println(err)
	}

	return money.Money(val), nil
}
