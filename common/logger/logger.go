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
	AsXml() ([]byte, error)
}

const (
	FUserCommandLog = "UserCommandLog"
	FQuoteServerLog = "QuoteServerLog"
	FAccountTransactionLog = "AccountTransactionLog"
	FSystemEventLog = "SystemEventLog"
	FErrorEventLog = "ErrorEventLog"
	FDebugEventLog = "DebugEventLog"
)


func (v UserCommandLog) AsXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v QuoteServerLog) AsXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v AccountTransactionLog) AsXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v SystemEventLog) AsXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v ErrorEventLog) AsXml() ([]byte, error) {
	return xml.Marshal(v)
}
func (v DebugEventLog) AsXml() ([]byte, error) {
  return xml.Marshal(v)
}

func (v ErrorEventLog) Error() string {
	return v.ErrorMessage
}

// Global logger worker
var logChan chan []byte
var dispatcher *gorpc.Dispatcher
var auditClient *gorpc.Client

func init() {

	gorpc.RegisterType(&UserCommandLog{})
	gorpc.RegisterType(&QuoteServerLog{})
	gorpc.RegisterType(&AccountTransactionLog{})
	gorpc.RegisterType(&SystemEventLog{})
	gorpc.RegisterType(&ErrorEventLog{})
	gorpc.RegisterType(&DebugEventLog{})

	dispatcher = gorpc.NewDispatcher()
	dispatcher.AddFunc(FUserCommandLog, func(v *UserCommandLog) error { return nil })
	dispatcher.AddFunc(FQuoteServerLog, func(v *QuoteServerLog) error { return nil })
	dispatcher.AddFunc(FAccountTransactionLog, func(v *AccountTransactionLog) error { return nil })
	dispatcher.AddFunc(FSystemEventLog , func(v *SystemEventLog) error { return nil })
	dispatcher.AddFunc(FErrorEventLog, func(v *ErrorEventLog) error { return nil })
	dispatcher.AddFunc(FDebugEventLog, func(v *DebugEventLog) error { return nil })

	auditClient = &gorpc.Client{
		Addr: fmt.Sprintf("%s:%d",
			config.GlobalConfig.AuditServer.Domain,
			config.GlobalConfig.AuditServer.Port),
	}

	fmt.Println("Connecting to audit server")
	auditClient.Start()
}

func Log(msg XmlLoggable) {
	client := dispatcher.NewFuncClient(auditClient)
	var err error
	switch v := msg.(type) {
	case UserCommandLog:
		_, err = client.Call(FUserCommandLog, v)
	case QuoteServerLog:
		_, err = client.Call(FQuoteServerLog, v)
	case AccountTransactionLog:
		_, err = client.Call(FAccountTransactionLog, v)
	case SystemEventLog:
		_, err = client.Call(FSystemEventLog, v)
	case ErrorEventLog:
		_, err = client.Call(FErrorEventLog, v)
	case DebugEventLog:
		_, err = client.Call(FDebugEventLog, v)
	}

	if err != nil {
		fmt.Println("Err logging, ", err)
	}
}



