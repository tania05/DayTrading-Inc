package main

import "fmt"
import "io/ioutil"
import "strings"
import "strconv"
import "net/http"
import "encoding/json"
import "bytes"

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

const lburl = "http://localhost:5555"
const post = "POST"
const put = "PUT"
const delete = "DELETE"
const get = "GET"

var loadBalanceMap map[string]string = make(map[string]string)


func postRequest(path string, url string, reqType string, payload interface{}) string{
  buff, _ := json.Marshal(payload)
  fmt.Println(url)
  req, e := http.NewRequest(strings.ToUpper(reqType), "http://"+url+path, bytes.NewBuffer(buff))
  if (e!= nil){
    panic(e)
  }
  req.Header.Add("Content-Type","application/json") 
  resp, e := http.DefaultClient.Do(req)

  if (e!= nil){
    panic(e)
  }
  
  defer resp.Body.Close()    
  bs, _ := ioutil.ReadAll(resp.Body)
  return fmt.Sprintln(string(bs))
}

func getRequest(path string, url string) string{ 
  req, _ := http.NewRequest("GET", "http://"+url+path, nil)
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
  return postRequest("/users", loadBalanceMap[command.UserId], post, command)
}

func (command QuoteCommand) request() string {
  // return fmt.Sprintln("QuoteCommand")  
  return postRequest("/stocks/"+command.StockSymbol+"/quote", loadBalanceMap[command.UserId], post, command)
}

func (command BuyCommand) request() string {
  // return fmt.Sprintln("BuyCommand")  
  return postRequest("/stocks/"+command.StockSymbol+"/buy", loadBalanceMap[command.UserId], post, command)
}

func (command CommitBuyCommand) request() string {
  // return fmt.Sprintln("CommitBuyCommand")  
  return postRequest("/stocks/buy", loadBalanceMap[command.UserId], put, command)
}

func (command CancelBuyCommand) request() string {
  // return fmt.Sprintln("CancelBuyCommand")
  return postRequest("/stocks/buy", loadBalanceMap[command.UserId], delete, command)
}

func (command SellCommand) request() string {
  // return fmt.Sprintln("SellCommand")
  return postRequest("/stocks/"+command.StockSymbol+"/sell", loadBalanceMap[command.UserId], post, command)
}

func (command CommitSellCommand) request() string {
  // return fmt.Sprintln("CommitSellCommand")
  return postRequest("/stocks/sell", loadBalanceMap[command.UserId], put, command)
}

func (command CancelSellCommand) request() string {
  // return fmt.Sprintln("CancelSellCommand")
  return postRequest("/stocks/sell", loadBalanceMap[command.UserId], delete, command)
}

func (command SetBuyAmountCommand) request() string {
  // return fmt.Sprintln("SetBuyAmountCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/buy", loadBalanceMap[command.UserId], post, command)
}

func (command CancelSetBuyCommand) request() string {
  // return fmt.Sprintln("CancelSetBuyCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/buy", loadBalanceMap[command.UserId], post, command)
}

func (command SetBuyTriggerCommand) request() string {
  // return fmt.Sprintln("SetBuyTriggerCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/buy", loadBalanceMap[command.UserId], post, command)
}

func (command SetSellAmountCommand) request() string {
  // return fmt.Sprintln("SetSellAmountCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/sell", loadBalanceMap[command.UserId], post, command)
}

func (command CancelSetSellCommand) request() string {
  // return fmt.Sprintln("CanceletSellCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, loadBalanceMap[command.UserId], command)
}

func (command SetSellTriggerCommand) request() string {
  // return fmt.Sprintln("SetSellTriggerCommand")
  return postRequest("/triggers/"+command.StockSymbol+"/sell", loadBalanceMap[command.UserId], post, command)
}

func (command DumplogCommand) request() string {
  // return fmt.Sprintln("DumplogCommand")
  return postRequest("/"+command.UserId+"/dump", loadBalanceMap[command.UserId], post, command)
}

func (command AdminDumblogCommand) request() string {
  // return fmt.Sprintln("AdminDumpLogCommand")
  return postRequest("/dump", lburl, post, command)
}

func (command DisplaySummary) request() string {
  // return fmt.Sprintln("DisplaySummaryCommand")
  return postRequest("/users/" + command.UserId + "/summary", loadBalanceMap[command.UserId], post, command)
}

func handleCommand(userMap map[string][]Command) {
  done := make(chan int, 250)
  for key, value := range userMap {
    go func(userId string, commands [] Command, done chan int){
      for _, n := range commands{
        fmt.Println(n.request())
      }
      done <- 0      
    }(key, value, done)
  }

  for range userMap {
    <-done
  }
}

func main() {
  b, err := ioutil.ReadFile("10userWorkLoad.txt")

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
    
    ip, ok := loadBalanceMap[userId]
    if (ok) {
      fmt.Println("Found ip: ", ip)
    } else {
      req, _ := http.NewRequest("GET", lburl+"/users", nil)
      req.Header.Add("Content-Type","application/json") 
      resp, e := http.DefaultClient.Do(req)

      if (e!= nil){
        panic(e)
      }
      defer resp.Body.Close()    
      bs, _ := ioutil.ReadAll(resp.Body)
      loadBalanceMap[userId] = string(bs)
      fmt.Println(string(bs))
    }

  
    

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
