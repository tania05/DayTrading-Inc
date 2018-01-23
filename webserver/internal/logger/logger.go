package logger

import (
	"encoding/xml"
	"os"
	"webserver/internal/money"
)

type commandType string
type action string

const (
	Add            commandType = "ADD"
	Quote          commandType = "QUOTE"
	Buy            commandType = "BUY"
	CommitBuy      commandType = "COMMIT_BUY"
	CancelBuy      commandType = "CANCEL_BUY"
	Sell           commandType = "SELL"
	CommitSell     commandType = "COMMIT_SELL"
	CancelSell     commandType = "CANCEL_SELL"
	SetBuyAmount   commandType = "SET_BUY_AMOUNT"
	SetBuyTrigger  commandType = "SET_BUY_TRIGGER"
	CancelSetBuy   commandType = "CANCEL_SET_BUY"
	SetSellAmount  commandType = "SET_SELL_AMOUNT"
	SetSellTrigger commandType = "SET_SELL_TRIGGER"
	CancelSetSell  commandType = "CANCEL_SET_SELL"
	DumpLog        commandType = "DUMPLOG"
	DisplaySummary commandType = "DISPLAY_SUMMARY"
)

const (
	AddAction    action = "add"
	RemoveAction action = "remove"
)

type UserCommandLog struct {
	XMLName        xml.Name    `xml:"userCommand"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int         `xml:"transactionNum"`
	Command        commandType `xml:"command"`
	Username       string      `xml:"username"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
}

type QuoteServerLog struct {
	XMLName         xml.Name    `xml:"quoteServer"`
	Timestamp       int64       `xml:"timestamp"`
	Server          string      `xml:"server"`
	TransactionNum  int         `xml:"transactionNum"`
	Price           money.Money `xml:"price"`
	StockSymbol     string      `xml:"stockSymbol"`
	Username        string      `xml:"username"`
	QuoteServerTime int64       `xml:"quoteServcerTime"`
	Cryptokey       string      `xml:"cryptokey"`
}

type AccountTransactionLog struct {
	XMLName        xml.Name    `xml:"accountTransaction"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int         `xml:"transactionNum"`
	Action         action      `xml:"action"`
	Username       string      `xml:"username,omitempty"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
}

type SystemEventLog struct {
	XMLName        xml.Name    `xml:"systemEvent"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int         `xml:"transactionNum"`
	Command        commandType `xml:"command"`
	Username       string      `xml:"username,omitempty"`
	StockSymbol    string      `xml:"stockSymbol,omitempty"`
	Filename       string      `xml:"filename,omitempty"`
	Funds          money.Money `xml:"funds,omitempty"`
}

type ErrorEventLog struct {
	XMLName        xml.Name    `xml:"errorEvent"`
	Timestamp      int64       `xml:"timestamp"`
	Server         string      `xml:"server"`
	TransactionNum int         `xml:"transactionNum"`
	Command        commandType `xml:"command"`
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
	TransactionNum int         `xml:"transactionNum"`
	Command        commandType `xml:"command"`
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
