package main

import (
	"github.com/valyala/gorpc"
	. "common/rpc/triggerstructs"
	"common/database"
	"common/config"
	"fmt"
	"log"
	"sync"
	"time"
	"strings"
	"common/quote"
	"common/context"
	"common/logger"
	"database/sql"
)

type trigger struct {
	executionPrice int
	amount int
	transactionNum int64
}

var buyTriggers map[string]*trigger
var sellTriggers map[string]*trigger
var mutex sync.Mutex

func executeBuyTrigger(userId string, stockSymbol string, price int, triggerId string, t *trigger, ctx *context.Context) {
	stockAmount := t.amount / price
	refund := t.amount - (stockAmount * price)

	db := database.GetDatabase(userId)
	tx, err := db.Begin()
	if err != nil {
		ctx.MakeError(err.Error())
		return
	}
	defer tx.Rollback()

	results, err := tx.Exec(`
		UPDATE users
			SET money = money + $2
			WHERE id = $1
	`, userId, refund)

	if err = ensureSingleRowAffected(err, ctx, results); err != nil {
		return
	}

	results, err = tx.Exec(`
			INSERT INTO stocks(user_id, stock_sym, amount)
				VALUES ($1, $2, $3)
				ON CONFLICT(user_id, stock_sym) DO UPDATE
					SET amount = stocks.amount + $3
		`, userId, stockSymbol, stockAmount)

	if err = ensureSingleRowAffected(err, ctx, results); err != nil {
		return
	}

	results, err = tx.Exec(`
		DELETE FROM triggers
			WHERE id = $1
	`, triggerId)

	if err = ensureSingleRowAffected(err, ctx, results); err != nil {
		delete(buyTriggers, triggerId)
		return
	}

	err = tx.Commit()
	if err != nil {
		mutex.Lock()
		defer mutex.Unlock()
		delete(buyTriggers, triggerId)
	}
}

func executeSellTrigger(userId string, stockSymbol string, price int, triggerId string, t *trigger, ctx *context.Context) {
	reservedStockAmount := int(t.amount / t.executionPrice)
	soldStockAmount := int(t.amount / price)
	recvMoney := price * soldStockAmount

	db := database.GetDatabase(userId)
	tx, err := db.Begin()
	if err != nil {
		ctx.MakeError(err.Error())
		return
	}
	defer tx.Rollback()

	results, err := tx.Exec(`
		UPDATE users
			SET money = money + $2
			WHERE id = $1
	`, userId, recvMoney)

	if err = ensureSingleRowAffected(err, ctx, results); err != nil {
		return
	}

	if soldStockAmount != reservedStockAmount {
		refundStockAmount := reservedStockAmount - soldStockAmount
		results, err = tx.Exec(`
			INSERT INTO stocks(user_id, stock_sym, amount)
				VALUES ($1, $2, $3)
				ON CONFLICT(user_id, stock_sym) DO UPDATE
					SET amount = stocks.amount + $3
			`, userId, stockSymbol, refundStockAmount)

		if err = ensureSingleRowAffected(err, ctx, results); err != nil {
			return
		}
	}

	results, err = tx.Exec(`
		DELETE FROM triggers
			WHERE id = $1
	`, triggerId)

	if err = ensureSingleRowAffected(err, ctx, results); err != nil {
		delete(sellTriggers, triggerId)
		return
	}

	err = tx.Commit()
	if err != nil {
		mutex.Lock()
		defer mutex.Unlock()
		delete(sellTriggers, triggerId)
	}
}

func ensureSingleRowAffected(err error, ctx *context.Context, results sql.Result) error {
	if err != nil {
		ctx.MakeError(err.Error())
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		ctx.MakeError(err.Error())
	}
	if rowsAffected != 1 {
		err = ctx.MakeError("Expected 1 row to be affected")
	}
	return err
}

func triggerLoop() {
	lastStart := time.Now()
	for {
		tts := lastStart.Add(60 * time.Second).Sub(time.Now())
		if tts < 0 {
			tts = 2000
		}
		time.Sleep(tts)
		lastStart = time.Now()
		fmt.Println("Beginning trigger execution check")

		wait := make(chan int)
		quoteWait := make(chan int, 6)
		quoteWait <- 1
		quoteWait <- 1
		quoteWait <- 1
		quoteWait <- 1
		total := 0

		mutex.Lock()
		for k, v := range buyTriggers {
			go func(k string, v *trigger) {
				stockSym := strings.Trim(k[:3], "!")
				userId := strings.Trim(k[3:67], "!")
				ctx := context.MakeSilentContext(v.transactionNum, userId, stockSym, logger.SetBuyTrigger)
				<-quoteWait
				price, err := quote.GetQuote(ctx)
				quoteWait <- 1
				if err != nil {
					return
				}
				if int(price) <= v.executionPrice {
					executeBuyTrigger(userId, stockSym, int(price), k, v, ctx)
				}
				wait <- 1
			}(k, v)
			total += 1
		}
		for k, v := range sellTriggers {
			go func(k string, v *trigger) {
				stockSym := strings.Trim(k[:3], "!")
				userId := strings.Trim(k[3:67], "!")
				ctx := context.MakeSilentContext(v.transactionNum, userId, stockSym, logger.SetSellTrigger)
				<-quoteWait
				price, err := quote.GetQuote(ctx)
				quoteWait <- 1
				if err != nil {
					return
				}
				if int(price) >= v.executionPrice {
					executeSellTrigger(userId, stockSym, int(price), k, v, ctx)
				}
				wait <- 1
			}(k, v)
			total += 1
		}
		mutex.Unlock()

		for i := 0; i < total; i++ {
			<-wait
		}
	}
}

func main() {
	buyTriggers = make(map[string]*trigger)
	sellTriggers = make(map[string]*trigger)

	d := gorpc.NewDispatcher()
	d.AddFunc(FSetBuyTriggerCommand, SetBuyTrigger)
	d.AddFunc(FSetBuyAmountCommand, SetBuyAmount)
	d.AddFunc(FSetSellTriggerCommand, SetSellTrigger)
	d.AddFunc(FSetSellAmountCommand, SetSellAmount)
	d.AddFunc(FCancelSetBuyCommand, CancelSetBuy)
	d.AddFunc(FCancelSetSellCommand, CancelSetSell)

	go triggerLoop()

	s := &gorpc.Server{
		Addr: fmt.Sprintf(":%d", config.GlobalConfig.Trigger.Port),
		Handler: d.NewHandlerFunc(),
	}

	if err := s.Serve(); err != nil {
		log.Fatalf("Can't start rpc server: %s", err)
		panic(err)
	}
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

func padStock(stock string) string {
	return padRight(stock, "!", 3)
}

func padUserId(userId string) string {
	return padRight(userId, "!", 64)
}

func getTriggerId(userId string, stockSym string, buySell bool) string {
	var buySellStr string
	if buySell {
		buySellStr = "t"
	} else {
		buySellStr = "f"
	}
	return padStock(stockSym) + padUserId(userId) + buySellStr
}

func setAmount(userId string, stockSymbol string, isBuy bool, amount int) error {
	db := database.GetDatabase(userId)
	triggerid := getTriggerId(userId, stockSymbol, isBuy)

	results, err := db.Exec(`
		INSERT INTO triggers(id, amount, is_buy, user_id, stock_sym)
			VALUES ($1, $2, $3, $4, $5)
	`, triggerid, amount, isBuy, userId, stockSymbol)

	if err != nil {
		fmt.Printf("setAmount err: %s\n", err)
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		fmt.Printf("rowsaffected %d: \n", rowsAffected)
	}


	return nil
}

func setTrigger(userId string, stockSymbol string, isBuy bool, executionPrice int, transactionNum int64) error {
	fmt.Println("Setting trigger")
	db := database.GetDatabase(userId)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	triggerId := getTriggerId(userId, stockSymbol, isBuy)
	fmt.Printf("trigger id %s\n", triggerId)

	row := tx.QueryRow(`
		UPDATE triggers
			SET execution_price = $1
			WHERE id = $2 AND execution_price IS NOT NULL 
			RETURNING amount
	`, executionPrice, triggerId)

	var amount int
	err = row.Scan(&amount)
	if err != nil {
		return err // TODO logger
	}

	if isBuy {
		_, err = tx.Exec(`
			UPDATE users
				SET money = money - $2
				WHERE id = $1
		`, userId, amount)
		if err != nil {
			return err
		}
	} else {
		stockAmount := int(amount/executionPrice)

		if stockAmount <= 0 {
			return fmt.Errorf("can not set trigger for 0 stocks\n")
		}

		_, err = tx.Exec(`
			UPDATE stocks
				SET amount = amount - $3
				WHERE user_id = $1 AND stock_sym = $2
		`, userId, stockSymbol, stockAmount)
		if err != nil {
			return err
		}
	}

	t := &trigger {
		amount: amount,
		executionPrice: executionPrice,
		transactionNum: transactionNum,
	}

	mutex.Lock()
	defer mutex.Unlock()
	if isBuy {
		buyTriggers[triggerId] = t
	} else {
		sellTriggers[triggerId] = t
	}

	tx.Commit()
	return nil
}

func cancelSet(userId string, stockSymbol string, isBuy bool) error {
	db := database.GetDatabase(userId)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	triggerId := getTriggerId(userId, stockSymbol, isBuy)

	var t *trigger
	mutex.Lock()
	if isBuy {
		if buyTriggers[triggerId] == nil { return fmt.Errorf("No trigger for userid/stocksymbol combination exists\n")}
		t = buyTriggers[triggerId]
	} else {
		if sellTriggers[triggerId] == nil { return fmt.Errorf("No trigger for userid/stocksymbol combination exists\n")}
		t = sellTriggers[triggerId]
	}
	mutex.Unlock()

	row := tx.QueryRow(`
		DELETE FROM triggers
			WHERE id = $1
			RETURNING (execution_price == NULL) as is_fully_set
	`, triggerId)

	var isFullySet bool
	err = row.Scan(&isFullySet)

	if err != nil {
		return err
	}

	var results sql.Result

	if isFullySet {
		if isBuy {
			results, err = tx.Exec(`
				UPDATE users
					SET money = money + $2
					WHERE id = $1
			`, userId, t.amount)
		} else {
			results, err = tx.Exec(`
				UPDATE stocks
				SET amount = amount + $3
				WHERE user_id = $1 AND stock_sym = $2
			`, userId, stockSymbol, int(t.amount/t.executionPrice))

		}
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("Unexpected number of rows affected during trigger refund")
	}

	tx.Commit()
	mutex.Lock()
	defer mutex.Unlock()
	if isBuy {
		delete(buyTriggers, triggerId)
	} else {
		delete(sellTriggers, triggerId)
	}


	return nil
}

func SetBuyAmount(request *SetBuyAmountCommand) error {
	return setAmount(request.UserId, request.StockSymbol, true, request.Amount)
}

func SetBuyTrigger(request *SetBuyTriggerCommand) error {
	return setTrigger(request.UserId, request.StockSymbol, true, request.ExecutionPrice, request.TransactionNum)
}

func SetSellAmount(request *SetSellAmountCommand) error {
	return setAmount(request.UserId, request.StockSymbol, false, request.Amount)
}

func SetSellTrigger(request *SetSellTriggerCommand) error {
	return setTrigger(request.UserId, request.StockSymbol, false, request.ExecutionPrice, request.TransactionNum)
}

func CancelSetBuy(request *CancelSetBuyCommand) error {
	return cancelSet(request.UserId, request.StockSymbol, true)
}

func CancelSetSell(request *CancelSetSellCommand) error {
	return cancelSet(request.UserId, request.StockSymbol, false)
}



//type trigger struct {
//  deposit        transaction.Holding
//  executionPrice money.Money
//  amount         money.Money
//}
//
//var triggerMutex sync.Mutex
//var buyTriggers map[string]map[string]*trigger
//var sellTriggers map[string]map[string]*trigger
//
//func init() {
//  buyTriggers = make(map[string]map[string]*trigger)
//  sellTriggers = make(map[string]map[string]*trigger)
//}
//
//func getUserMap(triggerMap map[string]map[string]*trigger, userId string) map[string]*trigger {
//  userMap := triggerMap[userId]
//  if userMap == nil {
//    triggerMap[userId] = make(map[string]*trigger)
//    return triggerMap[userId]
//  }
//
//  return userMap
//}
//
//func ensureTriggerDoesNotExist(ctx *context.Context, t *trigger) error {
//  if t != nil {
//    return ctx.MakeError("Trigger already exists")
//  }
//  return nil
//}
//
//func ensureTriggerExists(ctx *context.Context, t *trigger) error {
//  if t == nil {
//    return ctx.MakeError("Trigger does not exist")
//  }
//  return nil
//}
//
//func SetBuyAmount(ctx *context.Context, amount money.Money) error {
//  triggerMutex.Lock()
//  defer triggerMutex.Unlock()
//
//  userMap := getUserMap(buyTriggers, ctx.UserId)
//  t := userMap[ctx.StockSymbol]
//  if err := ensureTriggerDoesNotExist(ctx, t); err != nil {
//    return err
//  }
//
//  deposit, err := transaction.HoldMoney(ctx, amount)
//  if err != nil {
//    return err
//  }
//
//  userMap[ctx.StockSymbol] = &trigger {
//    deposit: deposit,
//    amount: amount,
//  }
//
//  return nil
//}
//
//func SetBuyTrigger(ctx *context.Context, executionPrice money.Money) error {
//  triggerMutex.Lock()
//  defer triggerMutex.Unlock()
//
//  userMap := getUserMap(buyTriggers, ctx.UserId)
//  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
//    return err
//  }
//
//  userMap[ctx.StockSymbol].executionPrice = executionPrice
//  return nil
//}
//
//func CancelSetBuy(ctx *context.Context) error {
//  triggerMutex.Lock()
//  defer triggerMutex.Unlock()
//
//  userMap := getUserMap(buyTriggers, ctx.UserId)
//  t := userMap[ctx.StockSymbol]
//  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
//    return err
//  }
//
//  transaction.Return(ctx, t.deposit)
//  delete(userMap, ctx.StockSymbol)
//
//  return nil
//}
//
//func SetSellAmount(ctx *context.Context, amount money.Money) error {
//  triggerMutex.Lock()
//  defer triggerMutex.Unlock()
//
//  userMap := getUserMap(sellTriggers, ctx.UserId)
//  t := userMap[ctx.StockSymbol]
//  if err := ensureTriggerDoesNotExist(ctx, t); err != nil {
//    return err
//  }
//
//  userMap[ctx.StockSymbol] = &trigger{
//    amount: amount,
//  }
//  return nil
//}
//
//func SetSellTrigger(ctx *context.Context, executionPrice money.Money) error {
//  triggerMutex.Lock()
//  defer triggerMutex.Unlock()
//
//  userMap := getUserMap(sellTriggers, ctx.UserId)
//  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
//    return err
//  }
//
//  requiredStocks := int(userMap[ctx.StockSymbol].amount) / int(executionPrice)
//  deposit, err := transaction.HoldStocks(ctx, requiredStocks)
//  if err != nil {
//    return err
//  }
//
//  userMap[ctx.StockSymbol].executionPrice = executionPrice
//  userMap[ctx.StockSymbol].deposit = deposit
//  return nil
//}
//
//func CancelSetSell(ctx *context.Context) error {
//  triggerMutex.Lock()
//  defer triggerMutex.Unlock()
//
//  userMap := getUserMap(sellTriggers, ctx.UserId)
//  t := userMap[ctx.StockSymbol]
//  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
//    return err
//  }
//
//  transaction.Return(ctx, t.deposit)
//  delete(userMap, ctx.StockSymbol)
//
//  return nil
//}
//
//func executeBuyTrigger(ctx *context.Context, t *trigger, quotePrice money.Money) error {
//  numStocks := int(t.amount) / int(quotePrice)
//  refundAmt := money.Money(int(t.amount) - (numStocks * int(quotePrice)))
//
//  reverseMoneyHold := transaction.MoneyHolding{UserId: ctx.UserId, Amount: refundAmt}
//  reverseStockHold := transaction.StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: numStocks}
//
//  err := transaction.Return(ctx, reverseMoneyHold, reverseStockHold)
//  return err
//}
//
//func executeSellTrigger(ctx *context.Context, t *trigger, quotePrice money.Money) error {
//  numStocks := int(t.amount) / int(quotePrice)
//  maxStocks := int(t.amount) / int(t.executionPrice)
//  refundStocks := maxStocks - numStocks
//  sellMoney := money.Money(numStocks * int(quotePrice))
//
//  reverseMoneyHold := transaction.MoneyHolding{UserId: ctx.UserId, Amount: sellMoney}
//  reverseStockHold := transaction.StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: refundStocks}
//
//  err := transaction.Return(ctx, reverseMoneyHold, reverseStockHold)
//  return err
//}

