package logger

import (
	"encoding/xml"
	"os"
	"webserver/internal/money"
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
	Username       string      `xml:"username"`
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

func init() {
	logChan = make(chan []byte, 32)
	go func() {
		f, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		for {
			bytes := <-logChan
			f.Write(bytes)
			f.Write([]byte("\n"))
		}
	}()
}

func Log(msg XmlLoggable) {
	bytes, err := msg.asXml()
	if err != nil { //TODO better error handling
		panic(err)
	}
	logChan <- bytes
}



