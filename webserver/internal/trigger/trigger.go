package trigger

import (
  "sync"
  "time"
  "webserver/internal/database"
  "webserver/internal/logger"
  "webserver/internal/money"
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

func ensureTriggerDoesNotExist(t *trigger, command logger.CommandType, userId string, stockSymbol string) error {
  if t != nil {
    err := logger.ErrorEventLog {
      Timestamp: time.Now().Unix(),
      Server: "ts0", //TODO
      TransactionNum: database.NewTransactionId(),
      Command: command,
      Username: userId,
      StockSymbol: stockSymbol,
      ErrorMessage: "Trigger already exists",
    }
    logger.Log(err)
    return err
  }
  return nil
}

func ensureTriggerExists(t *trigger, command logger.CommandType, userId string, stockSymbol string) error {
  if t == nil {
    err := logger.ErrorEventLog {
      Timestamp: time.Now().Unix(),
      Server: "ts0", //TODO
      TransactionNum: database.NewTransactionId(),
      Command: command,
      Username: userId,
      StockSymbol: stockSymbol,
      ErrorMessage: "Trigger does not exist",
    }
    logger.Log(err)
    return err
  }
  return nil
}

func ensureTriggerNotRunning(t *trigger, command logger.CommandType, userId string, stockSymbol string) error {
  if t == nil || t.executionPrice != 0 {
    err := logger.ErrorEventLog {
      Timestamp: time.Now().Unix(),
      Server: "ts0", //TODO
      TransactionNum: database.NewTransactionId(),
      Command: command,
      Username: userId,
      StockSymbol: stockSymbol,
      ErrorMessage: "Trigger for stock already running",
    }
    logger.Log(err)
    return err
  }
  return nil
}

func SetBuyAmount(userId string, stockSymbol string, amount money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(buyTriggers, userId)
  t := userMap[stockSymbol]
  if err := ensureTriggerDoesNotExist(t, logger.SetBuyAmount, userId, stockSymbol); err != nil {
    return err
  }

  deposit, err := database.HoldMoney(userId, amount)
  if err != nil {
    return err
  }

  userMap[stockSymbol] = &trigger {
    deposit: deposit,
    amount: amount,
  }

  return nil
}

func SetBuyTrigger(userId string, stockSymbol string, executionPrice money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(buyTriggers, userId)
  if err := ensureTriggerExists(userMap[stockSymbol], logger.SetBuyTrigger, userId, stockSymbol); err != nil {
    return err
  }

  userMap[stockSymbol].executionPrice = executionPrice
  return nil
}

func CancelSetBuy(userId string, stockSymbol string) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(buyTriggers, userId)
  t := userMap[stockSymbol]
  if err := ensureTriggerExists(userMap[stockSymbol], logger.CancelSetBuy, userId, stockSymbol); err != nil {
    return err
  }

  database.Return(t.deposit)
  delete(userMap, stockSymbol)

  return nil
}

func SetSellAmount(userId string, stockSymbol string, amount money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(sellTriggers, userId)
  t := userMap[stockSymbol]
  if err := ensureTriggerDoesNotExist(t, logger.SetSellAmount, userId, stockSymbol); err != nil {
    return err
  }

  userMap[stockSymbol] = &trigger{
    amount: amount,
  }
  return nil
}

func SetSellTrigger(userId string, stockSymbol string, executionPrice money.Money) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(sellTriggers, userId)
  if err := ensureTriggerExists(userMap[stockSymbol], logger.SetSellTrigger, userId, stockSymbol); err != nil {
    return err
  }

  requiredStocks := int(userMap[stockSymbol].amount) / int(executionPrice)
  deposit, err := database.HoldStocks(userId, stockSymbol, requiredStocks)
  if err != nil {
    return err
  }

  userMap[stockSymbol].executionPrice = executionPrice
  userMap[stockSymbol].deposit = deposit
  return nil
}

func CancelSetSell(userId string, stockSymbol string) error {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  userMap := getUserMap(sellTriggers, userId)
  t := userMap[stockSymbol]
  if err := ensureTriggerExists(userMap[stockSymbol], logger.CancelSetSell, userId, stockSymbol); err != nil {
    return err
  }

  database.Return(t.deposit)
  delete(userMap, stockSymbol)

  return nil

}

func executeBuyTrigger(t *trigger, quotePrice money.Money, userId string, stockSymbol string) error {
  numStocks := int(t.amount) / int(quotePrice)
  refundAmt := money.Money(int(t.amount) - (numStocks * int(quotePrice)))

  reverseMoneyHold := database.MoneyHolding{UserId: userId, Amount: refundAmt}
  reverseStockHold := database.StockHolding{UserId: userId, StockSymbol: stockSymbol, Amount: numStocks}

  err := database.Return(reverseMoneyHold, reverseStockHold)
  return err
}

func executeSellTrigger(t *trigger, quotePrice money.Money, userId string, stockSymbol string) error {
  numStocks := int(t.amount) / int(quotePrice)
  maxStocks := int(t.amount) / int(t.executionPrice)
  refundStocks := maxStocks - numStocks
  sellMoney := money.Money(numStocks * int(quotePrice))

  reverseMoneyHold := database.MoneyHolding{UserId: userId, Amount: sellMoney}
  reverseStockHold := database.StockHolding{UserId: userId, StockSymbol: stockSymbol, Amount: refundStocks}

  err := database.Return(reverseMoneyHold, reverseStockHold)
  return err
}

func OnQuoteUpdate(stockSymbol string, quotePrice money.Money) {
  triggerMutex.Lock()
  defer triggerMutex.Unlock()

  for userId, userMap := range buyTriggers {
    for triggerSym, t := range userMap {
      if triggerSym == stockSymbol {
        if t.executionPrice != 0 && quotePrice < t.executionPrice {
          executeBuyTrigger(t, quotePrice, userId, triggerSym)
          delete(userMap, triggerSym)
        }
      }
    }
  }

  for userId, userMap := range sellTriggers {
    for triggerSym, t := range userMap {
      if triggerSym == stockSymbol {
        if t.executionPrice != 0 && quotePrice > t.executionPrice {
          executeSellTrigger(t, quotePrice, userId, triggerSym)
          delete(userMap, triggerSym)
        }
      }
    }
  }
}

