package main

import (
  "net/http"
  "io"
  "encoding/json"
  "github.com/gorilla/mux"
  "log"
  "webserver/internal/money"
  "webserver/internal/trigger"
  "webserver/internal/logger"
  "webserver/internal/context"
  "io/ioutil"
	"common/config"
	"strconv"
	"fmt"
	"strings"
	"bytes"
	"net"
)


type RegisterServerCommand struct {
	IP string
	Port int
}

type AddCommand struct {
  TransactionNum int64
  UserId string
  Amount int
}

type QuoteCommand struct {
  TransactionNum int64  
  UserId string
  StockSymbol string
 }

type BuyCommand struct {
  TransactionNum int64  
  UserId string
  StockSymbol string
  Amount int
}

type CommitBuyCommand struct {
  TransactionNum int64  
  UserId string
}

type CancelBuyCommand struct {
  TransactionNum int64  
  UserId string
}

type SellCommand struct {
  TransactionNum int64  
  UserId string
  StockSymbol string
  Amount int
}

type CommitSellCommand struct {
  TransactionNum int64
  UserId string
}

type CancelSellCommand struct {
  TransactionNum int64
  UserId string
}

type SetBuyAmountCommand struct {
  TransactionNum int64
  UserId string
  StockSymbol string
  Amount int
}

type CancelSetBuyCommand struct {
  TransactionNum int64
  UserId string
  StockSymbol string
}

type SetBuyTriggerCommand struct {
  TransactionNum int64
  UserId string
  StockSymbol string
  Amount int
}

type SetSellAmountCommand struct {
  TransactionNum int64
  UserId string
  StockSymbol string
  Amount int
}

type CancelSetSellCommand struct {
  TransactionNum int64
  UserId string
  StockSymbol string
}

type SetSellTriggerCommand struct {
  TransactionNum int64
  UserId string
  StockSymbol string
  Amount int
}

type DumplogCommand struct {
  TransactionNum int64
  UserId string
  FileName string
}

type AdminDumplogCommand struct {
  TransactionNum int64
  FileName string
}

type DisplaySummary struct {
  TransactionNum int64
  UserId string
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello World!")
}

func PostAddHandler(w http.ResponseWriter, r *http.Request) {
	var payload AddCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.Add)
  err = addFunds(ctx, money.Money(payload.Amount))
	if err != nil {
		w.WriteHeader(400)
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, payload.UserId)
}

func GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var payload QuoteCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.Quote)
	getQuote(ctx)
	w.WriteHeader(http.StatusCreated)	
	io.WriteString(w, "Get Quote!\n")
	// io.WriteString(w, money)
}

func PostBuyHandler(w http.ResponseWriter, r *http.Request) {
	var payload BuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.Buy)
	transact(ctx, 1, money.Money(payload.Amount))
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Buy Command!")
	// fmt.Println(database.CheckFunds(payload.UserId))
	// fmt.Println(database.CheckStock(payload.UserId, payload.StockSymbol))
}

// commit
func PutBuyHandler(w http.ResponseWriter, r *http.Request) {
	var payload CommitBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.CommitBuy)	
	commitTransact(ctx, 1)
	io.WriteString(w, "Commit buy!")
	w.WriteHeader(http.StatusOK)	
	// fmt.Println(database.CheckFunds(payload.UserId))
}

// cancel
func DeleteBuyHandler(w http.ResponseWriter, r *http.Request) {
	var payload CancelBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.CancelBuy)		
	cancelTransact(ctx,1)
	io.WriteString(w, "Cancel Buy!")
	// fmt.Println(database.CheckFunds(payload.UserId))
	w.WriteHeader(http.StatusOK)	
}

func PostSellHandler(w http.ResponseWriter, r *http.Request) {
	var payload SellCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.Sell)			
	transact(ctx, 0, money.Money(payload.Amount))
	io.WriteString(w, "Sell!")
	w.WriteHeader(http.StatusOK)	
	// fmt.Println(database.CheckFunds(payload.UserId))
	// fmt.Println(database.CheckStock(payload.UserId, payload.StockSymbol))
}

// commit
func PutSellHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload CommitSellCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.CommitSell)	
	commitTransact(ctx, 0)
	io.WriteString(w, "Commit Sell!")
	// fmt.Println(database.CheckFunds(payload.UserId))
}

// cancel
func DeleteSellHandler(w http.ResponseWriter, r *http.Request) {
	var payload CancelSellCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.CancelSell)		
	cancelTransact(ctx,0)
	io.WriteString(w, "Cancel Sell!")
	w.WriteHeader(http.StatusOK)
	// fmt.Println(database.CheckFunds(payload.UserId))
}

//SetBuyAmount
func PostBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var payload SetBuyAmountCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetBuyAmount)			
	trigger.SetBuyAmount(ctx, money.Money(payload.Amount))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Set Buy Amount!")
}

//SetBuyTrigger
func PutBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var payload SetBuyTriggerCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetBuyTrigger)			
	trigger.SetBuyTrigger(ctx, money.Money(payload.Amount))
	w.WriteHeader(http.StatusOK)	
	io.WriteString(w, "Set Buy Trigger!")
}

func DeleteBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var payload CancelSetBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.CancelSetBuy)	
	trigger.CancelSetBuy(ctx)
	w.WriteHeader(http.StatusOK)	
	io.WriteString(w, "Cancel Buy Trigger!")
}

//SetSellAmount
func PostSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	var payload SetSellAmountCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetSellAmount)		
	trigger.SetSellAmount(ctx, money.Money(payload.Amount))
	io.WriteString(w, "Set sell amount!")
}

//SetSellTrigger
func PutSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	var payload SetSellTriggerCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetSellTrigger)		
	trigger.SetSellTrigger(ctx, money.Money(payload.Amount))
	io.WriteString(w, "Set Sell Trigger!")
}

func DeleteSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload CancelSetSellCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.CancelSetSell)		
	trigger.CancelSetSell(ctx)
	io.WriteString(w, "Cancel sell trigger!")
}

func PostAdminDumpLogHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Admin Dump logger")

	var payload AdminDumplogCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	context.MakeContext(payload.TransactionNum, "", "", logger.DumpLog)
}

func PostDumpLogHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "User Dump logger!")

	var payload DumplogCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.DumpLog)
}

func GetDisplaySummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	var payload DisplaySummary
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.DisplaySummary)

}

func postRequest(path string, reqType string, payload interface{}) string{
  buff, _ := json.Marshal(payload)
  req, _ := http.NewRequest(strings.ToUpper(reqType), path, bytes.NewBuffer(buff))
  req.Header.Add("Content-Type","application/json") 
  resp, e := http.DefaultClient.Do(req)

  if (e!= nil){
    panic(e)
  }
  
  defer resp.Body.Close()    
  bs, _ := ioutil.ReadAll(resp.Body)
  return string(bs)
}

func RegisterServer(port int) string{
	conn, err := net.Dial("udp", "1.2.3.4:80") //dummy connect, ill explain later.
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    payload := RegisterServerCommand{IP: localAddr.IP.String(), Port: port}
    return postRequest("http://" +
    	config.GlobalConfig.LoadBalancer.Domain + ":" + strconv.Itoa(config.GlobalConfig.LoadBalancer.Port) +
		"/register", "POST", payload)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.Path("/users").Methods("POST").HandlerFunc(PostAddHandler);
	r.Path("/stocks/{stockSym}/quote").Methods("POST").HandlerFunc(GetQuoteHandler);
	r.Path("/stocks/{stockSym}/buy").Methods("POST").HandlerFunc(PostBuyHandler);
	r.Path("/stocks/buy").Methods("PUT").HandlerFunc(PutBuyHandler);
	r.Path("/stocks/buy").Methods("DELETE").HandlerFunc(DeleteBuyHandler);
	r.Path("/stocks/{stockSym}/sell").Methods("POST").HandlerFunc(PostSellHandler);
	r.Path("/stocks/sell").Methods("PUT").HandlerFunc(PutSellHandler);
	r.Path("/stocks/sell").Methods("DELETE").HandlerFunc(DeleteSellHandler);
	r.Path("/triggers/{stockSym}/buy").Methods("POST").HandlerFunc(PostBuyTriggerHandler);
	r.Path("/triggers/{stockSym}/buy").Methods("PUT").HandlerFunc(PutBuyTriggerHandler);
	r.Path("/triggers/{stockSym}/buy").Methods("DELETE").HandlerFunc(DeleteBuyTriggerHandler);
	r.Path("/triggers/{stockSym}/sell").Methods("POST").HandlerFunc(PostSellTriggerHandler);
	r.Path("/triggers/{stockSym}/sell").Methods("PUT").HandlerFunc(PutSellTriggerHandler);
	r.Path("/triggers/{stockSym}/sell").Methods("DELETE").HandlerFunc(DeleteSellTriggerHandler);

	r.Path("/users/dump").Methods("POST").HandlerFunc(PostDumpLogHandler);
	r.Path("/dump").Methods("POST").HandlerFunc(PostAdminDumpLogHandler);
	r.Path("/users/{userId}/summary").Methods("POST").HandlerFunc(GetDisplaySummaryHandler);

	http.Handle("/", r)

	port := config.GlobalConfig.WebServer.Port
	addr := ":" + strconv.Itoa(port)
	fmt.Println(RegisterServer(port))
	log.Fatal(http.ListenAndServe(addr, nil))
}
