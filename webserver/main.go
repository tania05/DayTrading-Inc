package main

import (
  "net/http"
  "io"
  "encoding/json"
  "github.com/gorilla/mux"
  "fmt"
  "log"
  "webserver/internal/money"
  "webserver/internal/database"
  "webserver/internal/trigger"
  "io/ioutil"
)

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

type AdminDumblogCommand struct {
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
	w.WriteHeader(http.StatusOK)
	var payload AddCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	addFunds(payload.UserId, money.Money(payload.Amount))
	fmt.Println(database.CheckFunds(payload.UserId))
	io.WriteString(w, payload.UserId)
}

func GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Get Quote!")
}

func PostBuyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	var payload BuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	transact(1, payload.UserId, money.Money(payload.Amount), payload.StockSymbol)
	io.WriteString(w, "Buy Command!")
	fmt.Println(database.CheckFunds(payload.UserId))
	fmt.Println(database.CheckStock(payload.UserId, payload.StockSymbol))
}

// commit
func PutBuyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload CommitBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	commitTransact(1, payload.UserId)
	io.WriteString(w, "Commit buy!")
	fmt.Println(database.CheckFunds(payload.UserId))
}

// cancel
func DeleteBuyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload CancelBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	commitTransact(1, payload.UserId)
	io.WriteString(w, "Cancel Buy!")
	fmt.Println(database.CheckFunds(payload.UserId))
}

func PostSellHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload SellCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	transact(0, payload.UserId, money.Money(payload.Amount), payload.StockSymbol)
	io.WriteString(w, "Sell!")
	fmt.Println(database.CheckFunds(payload.UserId))
	fmt.Println(database.CheckStock(payload.UserId, payload.StockSymbol))
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
	commitTransact(0, payload.UserId)
	io.WriteString(w, "Commit Sell!")
	fmt.Println(database.CheckFunds(payload.UserId))
}

// cancel
func DeleteSellHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload CancelSellCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	commitTransact(0, payload.UserId)
	io.WriteString(w, "Cancel Sell!")
	fmt.Println(database.CheckFunds(payload.UserId))
}

//SetBuyAmount
func PostBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	var payload SetBuyAmountCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	trigger.SetBuyAmount(payload.UserId, payload.StockSymbol, money.Money(payload.Amount))
	io.WriteString(w, "Set Buy Amount!")
}

//SetBuyTrigger
func PutBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload SetBuyTriggerCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	trigger.SetBuyTrigger(payload.UserId, payload.StockSymbol, money.Money(payload.Amount))
	io.WriteString(w, "Set Buy Trigger!")
}

func DeleteBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var payload CancelSetBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	trigger.CancelSetBuy(payload.UserId, payload.StockSymbol)
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
	trigger.SetSellAmount(payload.UserId, payload.StockSymbol, money.Money(payload.Amount))
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
	trigger.SetSellTrigger(payload.UserId, payload.StockSymbol, money.Money(payload.Amount))
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
	trigger.CancelSetSell(payload.UserId, payload.StockSymbol)
	io.WriteString(w, "Cancel sell trigger!")
}

func PostAdminDumpLogHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Admin Dump logger")
}

func PostDumpLogHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "User Dump logger!")
}

func GetDisplaySummaryHandler(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)

	ret := map[string]interface{}{
		"total_buys":  3,
		"total_sells": 4,
		"user_id":     urlVars["userId"],
	}

	bytes, err := json.Marshal(ret) //TODO, we should be checking the error, but ehh... its fine, probably
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", HelloHandler)
  r.Path("/users").Methods("POST").HandlerFunc(PostAddHandler);
  r.Path("/stocks/{stockSym}").Methods("GET").HandlerFunc(GetQuoteHandler);
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
  r.Path("/users/{userId}/summary").Methods("GET").HandlerFunc(GetDisplaySummaryHandler);

  http.Handle("/", r)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
