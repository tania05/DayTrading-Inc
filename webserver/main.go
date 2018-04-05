package main

import (
  "net/http"
  "io"
  "encoding/json"
  "github.com/gorilla/mux"
  "log"
  "common/money"
  "common/logger"
  "common/context"
  "io/ioutil"
	"common/config"
	"strconv"
	"fmt"
	"strings"
	"bytes"
	"net"
	"github.com/valyala/gorpc"
	"os"
	"common/quote"
	"common/rpc/triggerstructs"
	"common/hashing"
	"webserver/internal/transaction"
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

var triggerServers []*gorpc.Client
var dispatcher *gorpc.Dispatcher
func setupTriggerRpcs() {
	triggerCountStr := os.Getenv("TRIGGER_COUNT")
	triggerCount, err := strconv.Atoi(triggerCountStr)
	if err != nil {
		panic(err)
	}

	gorpc.RegisterType(&logger.ErrorEventLog{})
	gorpc.RegisterType(&triggerstructs.SetBuyAmountCommand{})
	gorpc.RegisterType(&triggerstructs.SetBuyTriggerCommand{})
	gorpc.RegisterType(&triggerstructs.SetSellAmountCommand{})
	gorpc.RegisterType(&triggerstructs.SetSellTriggerCommand{})
	gorpc.RegisterType(&triggerstructs.CancelSetBuyCommand{})
	gorpc.RegisterType(&triggerstructs.CancelSetSellCommand{})

	dispatcher = gorpc.NewDispatcher()
	dispatcher.AddFunc(triggerstructs.FSetBuyAmountCommand, func(v *triggerstructs.SetBuyAmountCommand) error { return nil })
	dispatcher.AddFunc(triggerstructs.FSetBuyTriggerCommand, func(v *triggerstructs.SetBuyTriggerCommand) error { return nil })
	dispatcher.AddFunc(triggerstructs.FSetSellAmountCommand, func(v *triggerstructs.SetSellAmountCommand) error { return nil })
	dispatcher.AddFunc(triggerstructs.FSetSellTriggerCommand, func(v *triggerstructs.SetSellTriggerCommand) error { return nil })
	dispatcher.AddFunc(triggerstructs.FCancelSetBuyCommand, func(v *triggerstructs.CancelSetBuyCommand) error { return nil })
	dispatcher.AddFunc(triggerstructs.FCancelSetSellCommand, func(v *triggerstructs.CancelSetSellCommand) error { return nil })

	triggerServers = make([]*gorpc.Client, triggerCount)
	for i := 0; i < triggerCount; i++ {
		triggerServers[i] = &gorpc.Client{
			Addr: fmt.Sprintf("%s%d:%d",
				config.GlobalConfig.Trigger.Domain, i,
				config.GlobalConfig.Trigger.Port),
		}

		triggerServers[i].Start()
	}
}

func getTriggerClient(userId string) *gorpc.DispatcherClient {
	hashVal := hashing.ModuloHash(userId, len(triggerServers))
	return dispatcher.NewFuncClient(triggerServers[hashVal])
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
	newBalance, err := addFunds(ctx, money.Money(payload.Amount))
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, fmt.Sprintf("Succcess, your new balance is %d.\n", newBalance))
	}
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
	price, err := quote.GetQuote(ctx)

	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, fmt.Sprintf("Price of %s: %d\n", payload.StockSymbol, int(price)))
	}

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
	price, amount, err := transact(ctx, 1, money.Money(payload.Amount))
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("Started transaction - buy %d of %s for %d\n", amount, payload.StockSymbol, int(price)))
	}
	// fmt.Println(transaction.CheckFunds(payload.UserId))
	// fmt.Println(transaction.CheckStock(payload.UserId, payload.StockSymbol))
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
	err = commitTransact(ctx, 1)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Successfully committed last buy transaction")
	}
	// fmt.Println(transaction.CheckFunds(payload.UserId))
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
	err = cancelTransact(ctx,1)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Successfully canceled last buy transaction")
	}
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
	price, amount, err := transact(ctx, 0, money.Money(payload.Amount))
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("Started transaction - sell %d of %s for %d\n", amount, payload.StockSymbol, int(price)))
	}
	// fmt.Println(transaction.CheckFunds(payload.UserId))
	// fmt.Println(transaction.CheckStock(payload.UserId, payload.StockSymbol))
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
	err = commitTransact(ctx, 0)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Successfully committed last sell transaction")
	}
	// fmt.Println(transaction.CheckFunds(payload.UserId))
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
	err = cancelTransact(ctx,0)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Successfully canceled last sell transaction")
	}
	// fmt.Println(transaction.CheckFunds(payload.UserId))
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

	triggerClient := getTriggerClient(ctx.UserId)
	response, err := triggerClient.Call(triggerstructs.FSetBuyAmountCommand, triggerstructs.SetBuyAmountCommand{
		TransactionNum: ctx.TransactionNum,
		UserId: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
		Amount: payload.Amount,
	})

	if err != nil || response != nil{
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(200)
		io.WriteString(w, fmt.Sprintf("Began create buy trigger for %s for amount %d successfully", payload.StockSymbol, payload.Amount))
	}

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

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetBuyAmount)

	triggerClient := getTriggerClient(ctx.UserId)
	response, err := triggerClient.Call(triggerstructs.FSetBuyTriggerCommand, triggerstructs.SetBuyTriggerCommand{
		TransactionNum: ctx.TransactionNum,
		UserId: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
		ExecutionPrice: payload.Amount,
	})

	if err != nil || response != nil{
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprintf("Finished setting up  buy trigger for %s with execution price %d successfully", payload.StockSymbol, payload.Amount))
	}
}

func DeleteBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var payload CancelSetBuyCommand
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetBuyAmount)

	triggerClient := getTriggerClient(ctx.UserId)
	response, err := triggerClient.Call(triggerstructs.FCancelSetBuyCommand, triggerstructs.CancelSetBuyCommand{
		TransactionNum: ctx.TransactionNum,
		UserId: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
	})

	if err != nil || response != nil{
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprintf("Deleted or Canceled buy trigger successfully."))
	}
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

	triggerClient := getTriggerClient(ctx.UserId)
	response, err := triggerClient.Call(triggerstructs.FSetSellAmountCommand, triggerstructs.SetSellAmountCommand{
		TransactionNum: ctx.TransactionNum,
		UserId: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
		Amount: payload.Amount,
	})

	if err != nil || response != nil{
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(200)
		io.WriteString(w, fmt.Sprintf("Began create sell trigger for %s for amount %d successfully", payload.StockSymbol, payload.Amount))
	}
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

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetSellAmount)

	triggerClient := getTriggerClient(ctx.UserId)
	response, err := triggerClient.Call(triggerstructs.FSetSellTriggerCommand, triggerstructs.SetSellTriggerCommand{
		TransactionNum: ctx.TransactionNum,
		UserId: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
		ExecutionPrice: payload.Amount,
	})

	if err != nil || response != nil{
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprintf("Finished setting up sell trigger for %s with execution price %d successfully", payload.StockSymbol, payload.Amount))
	}
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

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, payload.StockSymbol, logger.SetBuyAmount)

	triggerClient := getTriggerClient(ctx.UserId)
	response, err := triggerClient.Call(triggerstructs.FCancelSetSellCommand, triggerstructs.CancelSetSellCommand{
		TransactionNum: ctx.TransactionNum,
		UserId: ctx.UserId,
		StockSymbol: ctx.StockSymbol,
	})

	if err != nil || response != nil{
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprintf("Deleted or Canceled sell trigger successfully."))
	}
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
	var dump logger.DumpAdmin
	dump.TransactionNum = payload.TransactionNum
	dump.FileName = payload.FileName
	logger.AdminDumplog(dump)
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
	var payload DisplaySummary
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	ctx := context.MakeContext(payload.TransactionNum, payload.UserId, "", logger.DisplaySummary)
	summary, err := transaction.Summary(ctx, payload.UserId)

	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, summary)
	}


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
	setupTriggerRpcs()

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

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/", r)

	port := config.GlobalConfig.WebServer.Port
	addr := ":" + strconv.Itoa(port)

	fmt.Printf("Listening on %d\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
