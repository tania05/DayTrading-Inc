package context

import (
  "webserver/internal/money"
  "webserver/internal/logger"
  "time"
)

type Context struct {
  TransactionNum int64
  UserId         string
  Command        logger.CommandType
  Funds          money.Money
  StockSymbol    string
}

func MakeContext(transactionNum int64, userId string, stockSymbol string, command logger.CommandType) *Context {
  return &Context{TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol, Command: command}
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
