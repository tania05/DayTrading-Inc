package main

import (
  "net/http"
  "io"
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
)
func HelloHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "Hello World!")
}

func PostAddHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func GetQuoteHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
}

func PostBuyHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
}

// commit
func PutBuyHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

// cancel
func DeleteBuyHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func PostSellHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
}

// commit
func PutSellHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

// cancel
func DeleteSellHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

//SetBuyAmount
func PostBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
}

//SetBuyTrigger
func PutBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func DeleteBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

//SetSellAmount
func PostSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
}

//SetSellTrigger
func PutSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}


func DeleteSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func PostDumpLogHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
}

type SampleSummaryResult struct{
  totalBuys int
  totalSells int
  userId string
}

func GetDisplaySummaryHandler(w http.ResponseWriter, r *http.Request) {
  urlVars := mux.Vars(r)

  ret := map[string]interface{}{
    "total_buys": 3,
    "total_sells": 4,
    "user_id": urlVars["userId"],
  }

  bytes, _ := json.Marshal(ret) //TODO, we should be checking the error, but ehh... its fine, probably
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
  r.Path("/stocks/{stockSym}/buy").Methods("PUT").HandlerFunc(PutBuyHandler);
  r.Path("/stocks/{stockSym}/buy").Methods("DELETE").HandlerFunc(DeleteBuyHandler);
  r.Path("/stocks/{stockSym}/sell").Methods("POST").HandlerFunc(PostSellHandler);
  r.Path("/stocks/{stockSym}/sell").Methods("PUT").HandlerFunc(PutSellHandler);
  r.Path("/stocks/{stockSym}/sell").Methods("DELETE").HandlerFunc(DeleteSellHandler);
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
