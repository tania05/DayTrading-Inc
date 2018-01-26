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
  "webserver/internal/database"
)
func HelloHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Hello World!")
}

func PostAddHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Add Post!")  
}

func GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Get Quote!") 
}

func PostBuyHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Buy!") 
}

// commit
func PutBuyHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Commit buy!")   
}

// cancel
func DeleteBuyHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Cancel Buy!")     
}

func PostSellHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  io.WriteString(w, "Sell!")
}

// commit
func PutSellHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Commit Sell!")
}

// cancel
func DeleteSellHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Cancel Sell!")
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

  fmt.Println(database.AddFunds("abcdef", 100))
  t, _ := database.AllocateFunds("abcdef", 10, "Q", 4)
  fmt.Println(t)
  fmt.Println(database.CheckFunds("abcdef"))
  fmt.Println(database.CheckStock("abcdef", "Q"))
  database.Commit(t)
  fmt.Println(database.CheckFunds("abcdef"))
  fmt.Println(database.CheckStock("abcdef", "Q"))

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

