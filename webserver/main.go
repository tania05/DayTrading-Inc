package main

import (
  "net/http"
  "io"
  "encoding/json"
  "encoding/xml"
  "github.com/gorilla/mux"
  "fmt"
  "log"
  "webserver/internal/logger"
  "webserver/internal/money"
  "webserver/internal/database"
  // "bytes"
  "io/ioutil"
)

type AddCommand struct {
  UserId string
  Amount int
}

type QuoteCommand struct {
  UserId string
  StockSymbol string
 }

type BuyCommand struct {
  UserId string
  StockSymbol string
  Amount int
}

type CommitBuyCommand struct {
  UserId string
}

type CancelBuyCommand struct {
  UserId string
}

type SellCommand struct {
  UserId string
  StockSymbol string
  Amount int
}

type CommitSellCommand struct {
  UserId string
}

type CancelSellCommand struct {
  UserId string
}

type SetBuyAmountCommand struct {
  UserId string
  StockSymbol string
  Amount int
}

type CancelSetBuyCommand struct {
  UserId string
  StockSymbol string
}

type SetBuyTriggerCommand struct {
  UserId string
  StockSymbol string
  Amount int
}

type SetSellAmountCommand struct {
  UserId string
  StockSymbol string
  Amount int
}

type CancelSetSellCommand struct {
  UserId string
  StockSymbol string
}

type SetSellTriggerCommand struct {
  UserId string
  StockSymbol string
  Amount int
}

type DumplogCommand struct {
  UserId string
  StockSymbol string
}

type DisplaySummary struct {
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
  transact(1,payload.UserId, money.Money(payload.Amount), payload.StockSymbol)  
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
  commitTransact(1,payload.UserId)
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
  commitTransact(1,payload.UserId)
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
  transact(0,payload.UserId, money.Money(payload.Amount), payload.StockSymbol)
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
  commitTransact(0,payload.UserId)
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
  commitTransact(0,payload.UserId)
  io.WriteString(w, "Cancel Sell!")
  fmt.Println(database.CheckFunds(payload.UserId))
}

//SetBuyAmount
func PostBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Set Buy Amount!")
}

//SetBuyTrigger
func PutBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Set Buy Trigger!")
}

func DeleteBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Cancel Buy Trigger!")
}

//SetSellAmount
func PostSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Set sell amount!")
}

//SetSellTrigger
func PutSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Set Sell Trigger!")
}

func DeleteSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Cancel sell trigger!")
}

func PostDumpLogHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Dump logger!")
}

func GetDisplaySummaryHandler(w http.ResponseWriter, r *http.Request) {
  urlVars := mux.Vars(r)

  ret := map[string]interface{}{
    "total_buys": 3,
    "total_sells": 4,
    "user_id": urlVars["userId"],
  }

  bytes, err := json.Marshal(ret) //TODO, we should be checking the error, but ehh... its fine, probably
  if(err != nil){
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(bytes)
}


func main() {
  l := logger.UserCommandLog{
    Command: logger.CommitBuy,
    TransactionNum: 2,
    Username: "abcdef",
    Timestamp: 1238099857,
    StockSymbol: "ABC",
    Funds: 100,
  }
  logger.Log(l)

  bytes, _ := xml.Marshal(l)
  fmt.Println(string(bytes))

  // fmt.Println(database.AddFunds("abcdef", 100))
  // t, _ := database.AllocateFunds("abcdef", 10, "Q", 4)
  // fmt.Println(t)
  // fmt.Println(database.CheckFunds("abcdef"))
  // fmt.Println(database.CheckStock("abcdef", "Q"))
  // database.Commit(t)
  // fmt.Println(database.CheckFunds("abcdef"))
  // fmt.Println(database.CheckStock("abcdef", "Q"))

  // addFunds("adad", 55000)
  // fmt.Println(database.CheckFunds("adad"))
  // transact(1, "adad", 50, "TTT")
  // commitTransact(1, "adad")
  // fmt.Println(database.CheckFunds("adad"))
  // fmt.Println(database.CheckStock("adad","TTT"))

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

  r.Path("/dump").Methods("POST").HandlerFunc(PostDumpLogHandler);
  r.Path("/users/{userId}/summary").Methods("GET").HandlerFunc(GetDisplaySummaryHandler);


  http.Handle("/", r)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

