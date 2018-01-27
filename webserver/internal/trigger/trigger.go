package trigger

import (
  "sync"
  "webserver/internal/database"
  "webserver/internal/money"
  "webserver/internal/context"
)

type trigger struct {
  deposit database.Holding
  executionPrice money.Money
  amount money.Money
}

var triggerMutex sync.Mutex
var buyTriggers map[string]map[string]*trigger
var sellTriggers map[string]map[string]*trigger

func init() {
  buyTriggers = make(map[string]map[string]*trigger)
  sellTriggers = make(map[string]map[string]*trigger)
}

func getUserMap(triggerMap map[string]map[string]*trigger, userId string) map[string]*trigger {
  userMap := triggerMap[userId]
  if userMap == nil {
    triggerMap[userId] = make(map[string]*trigger)
    return triggerMap[userId]
  }

  return userMap
}

func ensureTriggerDoesNotExist(ctx *context.Context, t *trigger) error {
  if t != nil {
    return ctx.MakeError("Trigger already exists")
  }
  return nil
}

func ensureTriggerExists(ctx *context.Context, t *trigger) error {
  if t == nil {
    return ctx.MakeError("Trigger does not exist")
  }
  return nil
}

func ensureTriggerNotRunning(ctx *context.Context,t *trigger) error {
  if t == nil || t.executionPrice != 0 {
    return ctx.MakeError("Trigger for stock already running")
  }
  return nil
}

func SetBuyAmount(ctx *context.Context, amount money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(buyTriggers, ctx.UserId)
  t := userMap[ctx.StockSymbol]
  if err := ensureTriggerDoesNotExist(ctx, t); err != nil {
    return err
  }

  deposit, err := database.HoldMoney(ctx, amount)
  if err != nil {
    return err
  }

  userMap[ctx.StockSymbol] = &trigger {
    deposit: deposit,
    amount: amount,
  }

  return nil
}

func SetBuyTrigger(ctx *context.Context, executionPrice money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(buyTriggers, ctx.UserId)
  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
    return err
  }

  userMap[ctx.StockSymbol].executionPrice = executionPrice
  return nil
}

func CancelSetBuy(ctx *context.Context) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(buyTriggers, ctx.UserId)
  t := userMap[ctx.StockSymbol]
  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
    return err
  }

  database.Return(ctx, t.deposit)
  delete(userMap, ctx.StockSymbol)

  return nil
}

func SetSellAmount(ctx *context.Context, amount money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(sellTriggers, ctx.UserId)
  t := userMap[ctx.StockSymbol]
  if err := ensureTriggerDoesNotExist(ctx, t); err != nil {
    return err
  }

  userMap[ctx.StockSymbol] = &trigger{
    amount: amount,
  }
  return nil
}

func SetSellTrigger(ctx *context.Context, executionPrice money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(sellTriggers, ctx.UserId)
  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
    return err
  }

  requiredStocks := int(userMap[ctx.StockSymbol].amount) / int(executionPrice)
  deposit, err := database.HoldStocks(ctx, requiredStocks)
  if err != nil {
    return err
  }

  userMap[ctx.StockSymbol].executionPrice = executionPrice
  userMap[ctx.StockSymbol].deposit = deposit
  return nil
}

func CancelSetSell(ctx *context.Context) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(sellTriggers, ctx.UserId)
  t := userMap[ctx.StockSymbol]
  if err := ensureTriggerExists(ctx, userMap[ctx.StockSymbol]); err != nil {
    return err
  }

  database.Return(ctx, t.deposit)
  delete(userMap, ctx.StockSymbol)

  return nil
}

func executeBuyTrigger(ctx *context.Context, t *trigger, quotePrice money.Money) error {
  numStocks := int(t.amount) / int(quotePrice)
  refundAmt := money.Money(int(t.amount) - (numStocks * int(quotePrice)))

  reverseMoneyHold := database.MoneyHolding{UserId: ctx.UserId, Amount: refundAmt}
  reverseStockHold := database.StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: numStocks}

  err := database.Return(ctx, reverseMoneyHold, reverseStockHold)
  return err
}

func executeSellTrigger(ctx *context.Context, t *trigger, quotePrice money.Money) error {
  numStocks := int(t.amount) / int(quotePrice)
  maxStocks := int(t.amount) / int(t.executionPrice)
  refundStocks := maxStocks - numStocks
  sellMoney := money.Money(numStocks * int(quotePrice))

  reverseMoneyHold := database.MoneyHolding{UserId: ctx.UserId, Amount: sellMoney}
  reverseStockHold := database.StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: refundStocks}

  err := database.Return(ctx, reverseMoneyHold, reverseStockHold)
  return err
}

func OnQuoteUpdate(ctx *context.Context, quotePrice money.Money) {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  oldCtxUserId := ctx.UserId
  for userId, userMap := range buyTriggers {
    ctx.UserId = userId
    for triggerSym, t := range userMap {
      if triggerSym == ctx.StockSymbol {
        if t.executionPrice != 0 && quotePrice < t.executionPrice {
          executeBuyTrigger(ctx, t, quotePrice)
          delete(userMap, triggerSym)
        }
      }
    }
    ctx.UserId = oldCtxUserId
  }

  for userId, userMap := range sellTriggers {
    ctx.UserId = userId
    for triggerSym, t := range userMap {
      if triggerSym == ctx.StockSymbol {
        if t.executionPrice != 0 && quotePrice > t.executionPrice {
          executeSellTrigger(ctx, t, quotePrice)
          delete(userMap, triggerSym)
        }
      }
    }
    ctx.UserId = oldCtxUserId
  }
}

