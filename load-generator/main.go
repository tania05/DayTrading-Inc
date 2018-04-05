package main

import "fmt"
import "io/ioutil"
import "strings"
import "strconv"
import "net/http"
import "encoding/json"
import (
  "bytes"
  "os"
  "time"
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

type Command interface {
  request() string
}

const url = "http://10.0.75.1:8000"
const post = "POST"
const put = "PUT"
const delete = "DELETE"
const get = "GET"

var errorCount int

func postRequest(path string, reqType string, payload interface{}) string{
  buff, _ := json.Marshal(payload)
  req, _ := http.NewRequest(strings.ToUpper(reqType), url+path, bytes.NewBuffer(buff))
  req.Header.Add("Content-Type","application/json") 
  resp, e := http.DefaultClient.Do(req)

  if (e!= nil){
    time.Sleep(2 * time.Second)
    fmt.Println("error: ", e)
    return postRequest(path, reqType, payload)
  } else {
    defer resp.Body.Close()
    bs, _ := ioutil.ReadAll(resp.Body)
    return fmt.Sprintln(string(bs))
  }
  
}

func getRequest(path string) string{ 
  req, _ := http.NewRequest("GET", url+path, nil)
  resp, e := http.DefaultClient.Do(req)
  if (e!= nil){
    panic(e)
  }
  
  defer resp.Body.Close()    
  bs, _ := ioutil.ReadAll(resp.Body)
  return fmt.Sprintln(string(bs))
}

func (command AddCommand) request() string {
  //return fmt.Sprintln("AddCommand")
  return postRequest("/users", post, command)
}

func (command QuoteCommand) request() string {
  // return fmt.Sprintln("QuoteCommand")  
  return postRequest("/stocks/"+command.StockSymbol+"/quote", post, command)
}

func (command BuyCommand) request() string {
  // return fmt.Sprintln("BuyCommand")  
  return postRequest("/stocks/"+command.StockSymbol+"/buy", post, command)
}

func (command CommitBuyCommand) request() string {
  // return fmt.Sprintln("CommitBuyCommand")  
  return postRequest("/stocks/buy", put, command)
}

func (command CancelBuyCommand) request() string {
  // return fmt.Sprintln("CancelBuyCommand")
  return postRequest("/stocks/buy", delete, command)
}

func (command SellCommand) request() string {
  // return fmt.Sprintln("SellCommand")
  return postRequest("/stocks/"+command.StockSymbol+"/sell", post, command)
}

func (command CommitSellCommand) request() string {
  // return fmt.Sprintln("CommitSellCommand")
  return postRequest("/stocks/sell", put, command)
}

func (command CancelSellCommand) request() string {
  // return fmt.Sprintln("CancelSellCommand")
  return postRequest("/stocks/sell", delete, command)
}

func (command SetBuyAmountCommand) request() string {
  // return fmt.Sprintln("SetBuyAmountCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/buy", post, command)
}

func (command CancelSetBuyCommand) request() string {
  // return fmt.Sprintln("CancelSetBuyCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/buy", post, command)
}

func (command SetBuyTriggerCommand) request() string {
  // return fmt.Sprintln("SetBuyTriggerCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/buy", post, command)
}

func (command SetSellAmountCommand) request() string {
  // return fmt.Sprintln("SetSellAmountCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, command)
}

func (command CancelSetSellCommand) request() string {
  // return fmt.Sprintln("CanceletSellCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, command)
}

func (command SetSellTriggerCommand) request() string {
  // return fmt.Sprintln("SetSellTriggerCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, command)
}

func (command DumplogCommand) request() string {
  // return fmt.Sprintln("DumplogCommand")
  return postRequest("/"+command.UserId+"/dump", post, command)
}

func (command AdminDumblogCommand) request() string {
  // return fmt.Sprintln("AdminDumpLogCommand")
  return postRequest("/dump", post, command)
}

func (command DisplaySummary) request() string {
  // return fmt.Sprintln("DisplaySummaryCommand")
  return postRequest("/users/" + command.UserId + "/summary", post, command)
}

func handleCommand(userMap map[string][]Command) {
  fmt.Println("Total users: ", len(userMap))
  gen_id_str := os.Getenv("GENERATOR_ID")
  total_gen_str := os.Getenv("TOTAL_GENERATORS")

  gen_id, err := strconv.Atoi(gen_id_str)
  if err != nil {
    panic(err)
  }

  total_gen, err := strconv.Atoi(total_gen_str)
  if err != nil {
    panic(err)
  }

  i := 0

  done := make(chan int, 250)
  count := 0
  for key, value := range userMap {
    i++
    if i % total_gen != gen_id {
      continue
    }
    count++
    go func(userId string, commands [] Command, done chan int){
      for k, n := range commands{
        fmt.Println(n.request())
        fmt.Println("[", k, "/", len(commands), "]")
      }
      done <- 0      
    }(key, value, done)
  }

  fmt.Println("count", count)
  for fin := 0; fin < count; fin++ {
    <-done
  }
}

func main() {
  b, err := ioutil.ReadFile("testDump.txt")

  if (err != nil) {
    panic(err)
  }

  lines := strings.Split(string(b), "\n")
  //// "string", [] if not empty
  userMap := make(map[string][]Command)
  for _, line := range lines {
    line = strings.Trim(line, "\n\r \t")
    parts := strings.Split(line, ",")
    if len(parts) < 2 {
      continue
    }
    // fmt.Println("************************************************************")
    // fmt.Println(parts[0])
    userId := parts[1]
    transCommand := strings.Split(parts[0], " ")
    transactionNum, _ := strconv.ParseInt(strings.Trim((transCommand[0]),"[]"),10,0)
    commandType := transCommand[1]
    // fmt.Println(transactionNum)
    // fmt.Println(commandType)
    // transactionNum, commandType := strings.Split(parts[0], " ")
    
    switch commandType{
    case "ADD":
      a64, err := strconv.ParseInt(strings.Replace(parts[2],".","",-1),10,0)
      if(err !=nil){
        panic(err)
      }
      var amount = int(a64)
      
      addCommand := AddCommand{TransactionNum: transactionNum, UserId: userId, Amount: amount}
      userMap[userId] = append(userMap[userId], addCommand)
    
    case "QUOTE":
      stockSymbol := parts[2]
      quoteCommand := QuoteCommand{TransactionNum: transactionNum,UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], quoteCommand)
    
    case "BUY":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil){
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      buyCommand := BuyCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], buyCommand)
     
    case "COMMIT_BUY":
      commitBuyCommand := CommitBuyCommand{ TransactionNum: transactionNum, UserId : userId }
      userMap[userId] = append(userMap[userId], commitBuyCommand)
    
    case "CANCEL_BUY":
      cancelBuyCommand := CancelBuyCommand{ TransactionNum: transactionNum, UserId : userId}
      userMap[userId] = append(userMap[userId], cancelBuyCommand)

    case "SELL":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil){
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      sellCommand := SellCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], sellCommand)
     
    case "COMMIT_SELL":
      commitSellCommand := CommitSellCommand{ TransactionNum: transactionNum, UserId : userId}
      userMap[userId] = append(userMap[userId], commitSellCommand)
    
    case "CANCEL_SELL":
      cancelSellCommand := CancelSellCommand{ TransactionNum: transactionNum, UserId : userId}
      userMap[userId] = append(userMap[userId], cancelSellCommand)
    
    case "SET_BUY_AMOUNT":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil) {
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      setBuyAmountCommand := SetBuyAmountCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], setBuyAmountCommand)
    
    case "CANCEL_SET_BUY":
      stockSymbol := parts[2]
      cancelSetBuyCommand := CancelSetBuyCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], cancelSetBuyCommand)
    
    case "SET_BUY_TRIGGER":
      stockSymbol := parts[2]
      setBuyTriggerCommand := SetBuyTriggerCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], setBuyTriggerCommand)
    
    case "SET_SELL_AMOUNT":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil) {
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      setSellAmountCommand := SetSellAmountCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], setSellAmountCommand)
    
    case "CANCEL_SET_SELL":
      stockSymbol := parts[2]
      cancelSetBuyCommand := CancelSetBuyCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], cancelSetBuyCommand)
    
    case "SET_SELL_TRIGGER":
      stockSymbol := parts[2]
      setSellTriggerCommand := SetSellTriggerCommand{ TransactionNum: transactionNum, UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], setSellTriggerCommand)
    
    case "DUMPLOG":
      if(len(parts) < 3) {
        fileName := parts[1]
        dumplogCommand := AdminDumblogCommand{ TransactionNum: transactionNum, FileName: fileName}
        userMap["\xFF\xFF"] = append(userMap["\xFF\xFF"], dumplogCommand)
      } else {
        fileName := parts[2]
        dumplogCommand := DumplogCommand{ TransactionNum: transactionNum, UserId: userId, FileName: fileName}
        userMap[userId] = append(userMap[userId], dumplogCommand)
      }
      
    case "DISPLAY_SUMMARY":
      displaySummary := DisplaySummary{ TransactionNum: transactionNum, UserId : userId}
      userMap[userId] = append(userMap[userId], displaySummary)
    }
  }
  handleCommand(userMap)
}
