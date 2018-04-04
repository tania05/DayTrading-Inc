package logger

import (
	"encoding/xml"
	"common/money"
	"github.com/valyala/gorpc"
	"fmt"
	"common/config"
)

type CommandType string
type Action string

const (
	Add            CommandType = "ADD"
	Quote          CommandType = "QUOTE"
	Buy            CommandType = "BUY"
	CommitBuy      CommandType = "COMMIT_BUY"
	CancelBuy      CommandType = "CANCEL_BUY"
	Sell           CommandType = "SELL"
	CommitSell     CommandType = "COMMIT_SELL"
	CancelSell     CommandType = "CANCEL_SELL"
	SetBuyAmount   CommandType = "SET_BUY_AMOUNT"
	SetBuyTrigger  CommandType = "SET_BUY_TRIGGER"
	CancelSetBuy   CommandType = "CANCEL_SET_BUY"
	SetSellAmount  CommandType = "SET_SELL_AMOUNT"
	SetSellTrigger CommandType = "SET_SELL_TRIGGER"
	CancelSetSell  CommandType = "CANCEL_SET_SELL"
	DumpLog        CommandType = "DUMPLOG"
	DisplaySummary CommandType = "DISPLAY_SUMMARY"
)

const (
	AddAction    Action = "add"
	RemoveAction Action = "remove"
)

type UserCommandLog struct {
	XMLName        xml.Name    `xml:"userCommand"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int64       `xml:"transactionNum"`
	Command        CommandType `xml:"command"`
	Username       string      `xml:"username,omitempty"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
}

type QuoteServerLog struct {
	XMLName         xml.Name    `xml:"quoteServer"`
	Timestamp       int64       `xml:"timestamp"`
	Server          string      `xml:"server"`
	TransactionNum  int64       `xml:"transactionNum"`
	Price           money.Money `xml:"price"`
	StockSymbol     string      `xml:"stockSymbol"`
	Username        string      `xml:"username"`
	QuoteServerTime int64       `xml:"quoteServerTime"`
	Cryptokey       string      `xml:"cryptokey"`
}

type AccountTransactionLog struct {
	XMLName        xml.Name    `xml:"accountTransaction"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int64       `xml:"transactionNum"`
	Action         Action      `xml:"action"`
	Username       string      `xml:"username"`
	Funds          money.Money `xml:"funds"`
}

type SystemEventLog struct {
	XMLName        xml.Name    `xml:"systemEvent"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int64       `xml:"transactionNum"`
	Command        CommandType `xml:"command"`
	Username       string      `xml:"username,omitempty"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
}

type ErrorEventLog struct {
	XMLName        xml.Name    `xml:"errorEvent"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int64       `xml:"transactionNum"`
	Command        CommandType `xml:"command"`
	Username       string      `xml:"username,omitempty"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
	ErrorMessage   string      `xml:"errorMessage,omitempty"`
}

type DebugEventLog struct {
	XMLName        xml.Name    `xml:"debugEvent"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int64       `xml:"transactionNum"`
	Command        CommandType `xml:"command"`
	Username       string      `xml:"username,omitempty"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
	DebugMessage   string      `xml:"debugMessage,omitempty"`
}

type XmlLoggable interface {
	asXml() ([]byte, error)
}

func (v UserCommandLog) asXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v QuoteServerLog) asXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v AccountTransactionLog) asXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v SystemEventLog) asXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v ErrorEventLog) asXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v DebugEventLog) asXml() ([]byte, error) {
  return xml.Marshal(v)
}

func (v ErrorEventLog) Error() string {
  bytes, internalErr := v.asXml()
  if internalErr != nil {
    return "Error creating error message"
  }

  return string(bytes)
}

// Global logger worker
var logChan chan []byte
var dispatcher *gorpc.Dispatcher
var auditClient *gorpc.Client

func init() {

	gorpc.RegisterType(UserCommandLog{})
	gorpc.RegisterType(QuoteServerLog{})
	gorpc.RegisterType(AccountTransactionLog{})
	gorpc.RegisterType(SystemEventLog{})
	gorpc.RegisterType(ErrorEventLog{})
	gorpc.RegisterType(DebugEventLog{})

	dispatcher = gorpc.NewDispatcher()
	dispatcher.AddFunc(FLog, func(v *XmlLoggable) error { return nil })

	auditClient := &gorpc.Client{
		Addr: fmt.Sprintf("%s:%d",
			config.GlobalConfig.Trigger.Domain,
			config.GlobalConfig.Trigger.Port),
	}

	auditClient.Start()
	}

const FLog = "Log"
func Log(msg XmlLoggable) {
	client := dispatcher.NewFuncClient(auditClient)
	client.Send(FLog, &msg)
}



