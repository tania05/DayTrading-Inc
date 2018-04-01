package context

import (
  "webserver/internal/money"
  "webserver/internal/logger"
  "time"
  "fmt"
)

type Context struct {
  TransactionNum int64
  UserId         string
  Command        logger.CommandType
  Funds          money.Money
  StockSymbol    string
}

func MakeContext(transactionNum int64, userId string, stockSymbol string, command logger.CommandType) *Context {
  ctx := Context{TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol, Command: command}
  logger.Log(logger.UserCommandLog {
    StockSymbol: ctx.StockSymbol,
    Username: ctx.UserId,
    TransactionNum: ctx.TransactionNum,
    Command: ctx.Command,
    Funds: ctx.Funds,
    Server: "ts0",
    Timestamp: time.Now().UnixNano() / 1e6,
  })

  return &ctx
}

func (context Context) MakeError(message string) logger.ErrorEventLog {
  err := logger.ErrorEventLog{
    Timestamp:      time.Now().UnixNano() / 1e6,
    Server:         "ts0",
    TransactionNum: context.TransactionNum,
    Command:        context.Command,
    ErrorMessage: message,
    Funds: context.Funds,
  }
  logger.Log(err)
  fmt.Printf("%s\n", err)
  return err
}

func (context Context) MakeAccountTransactionLog(action logger.Action) {
  logger.Log(logger.AccountTransactionLog{
    Timestamp: time.Now().UnixNano() / 1e6,
    Server: "ts0",
    TransactionNum: context.TransactionNum,
    Action: action,
    Username: context.UserId,
  })
}
