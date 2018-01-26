package main

import "fmt"
import "io/ioutil"
import "strings"
import "strconv"
import "net/http"
import "encoding/json"
import "bytes"

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

type Command interface {
  request() string
}

const url = "http://webserver:8080"
const post = "POST"
const put = "PUT"
const delete = "DELETE"
const get = "GET"

func postRequest(path string, reqType string, payload interface{}) string{
  buff, _ := json.Marshal(payload)
  req, _ := http.NewRequest(strings.ToUpper(reqType), url+path, bytes.NewBuffer(buff))
  resp, e := http.DefaultClient.Do(req)

  if (e!= nil){
    panic(e)
  }
  
  defer resp.Body.Close()    
  bs, _ := ioutil.ReadAll(resp.Body)
  return fmt.Sprintln(string(bs))
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
  return postRequest("/users", post, command)
}

func (command QuoteCommand) request() string {
  return getRequest("/users")
}

func (command BuyCommand) request() string {
  return postRequest("/stocks/"+command.StockSymbol+"/buy", post, command)
}

func (command CommitBuyCommand) request() string {
  return postRequest("/stocks/buy", put, command)
}

func (command CancelBuyCommand) request() string {
  return postRequest("/stocks/buy", delete, command)
}

func (command SellCommand) request() string {
  return postRequest("/stocks/"+command.StockSymbol+"/sell", post, command)
}

func (command CommitSellCommand) request() string {
  return postRequest("/stocks/sell", put, command)
}

func (command CancelSellCommand) request() string {
  return postRequest("/stocks/sell", delete, command)
}

func (command SetBuyAmountCommand) request() string {
  return postRequest("/triggers/"+command.StockSymbol+"/buy", post, command)
}

func (command CancelSetBuyCommand) request() string {
  return postRequest("/triggers/"+command.StockSymbol+"/buy", post, command)
}

func (command SetBuyTriggerCommand) request() string {
  return postRequest("/triggers/"+command.StockSymbol+"/buy", post, command)
}

func (command SetSellAmountCommand) request() string {
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, command)
}

func (command CancelSetSellCommand) request() string {
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, command)
}

func (command SetSellTriggerCommand) request() string {
  return postRequest("/triggers/"+command.StockSymbol+"/sell", post, command)
}

func (command DumplogCommand) request() string {
  return postRequest("/dump", post, command)
}

func (command DisplaySummary) request() string {
  return getRequest("/users/" + command.UserId + "/summary")
}

func handleCommand(userMap map[string][]Command) {
  for key, value := range userMap {
    fmt.Println(key)
    for _, n := range value{
      fmt.Println(n.request())
    }
  }
}

func main() {
  b, err := ioutil.ReadFile("1userWorkLoad.txt")

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
    
    userId := parts[1]

    commandType := strings.Split(parts[0], " ")[1]
    
    switch commandType{
    case "ADD":
      a64, err := strconv.ParseInt(strings.Replace(parts[2],".","",-1),10,0)
      if(err !=nil){
        panic(err)
      }
      var amount = int(a64)
      
      addCommand := AddCommand{UserId: userId, Amount: amount}
      userMap[userId] = append(userMap[userId], addCommand)
    
    case "QUOTE":
      stockSymbol := parts[2]
      quoteCommand := QuoteCommand{UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], quoteCommand)
    
    case "BUY":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil){
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      buyCommand := BuyCommand{UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], buyCommand)
     
    case "COMMIT_BUY":
      commitBuyCommand := CommitBuyCommand{UserId : userId}
      userMap[userId] = append(userMap[userId], commitBuyCommand)
    
    case "CANCEL_BUY":
      cancelBuyCommand := CancelBuyCommand{UserId : userId}
      userMap[userId] = append(userMap[userId], cancelBuyCommand)

    case "SELL":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil){
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      sellCommand := SellCommand{UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], sellCommand)
     
    case "COMMIT_SELL":
      commitSellCommand := CommitSellCommand{UserId : userId}
      userMap[userId] = append(userMap[userId], commitSellCommand)
    
    case "CANCEL_SELL":
      cancelSellCommand := CancelSellCommand{UserId : userId}
      userMap[userId] = append(userMap[userId], cancelSellCommand)
    
    case "SET_BUY_AMOUNT":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil) {
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      setBuyAmountCommand := SetBuyAmountCommand{UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], setBuyAmountCommand)
    
    case "CANCEL_SET_BUY":
      stockSymbol := parts[2]
      cancelSetBuyCommand := CancelSetBuyCommand{UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], cancelSetBuyCommand)
    
    case "SET_BUY_TRIGGER":
      stockSymbol := parts[2]
      setBuyTriggerCommand := SetBuyTriggerCommand{UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], setBuyTriggerCommand)
    
    case "SET_SELL_AMOUNT":
      a64, err := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      if(err !=nil) {
        panic(err)
      }
      var amount = int(a64)
      stockSymbol := parts[2]
      
      setSellAmountCommand := SetSellAmountCommand{UserId: userId, StockSymbol: stockSymbol, Amount: amount}
      userMap[userId] = append(userMap[userId], setSellAmountCommand)
    
    case "CANCEL_SET_SELL":
      stockSymbol := parts[2]
      cancelSetBuyCommand := CancelSetBuyCommand{UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], cancelSetBuyCommand)
    
    case "SET_SELL_TRIGGER":
      stockSymbol := parts[2]
      setSellTriggerCommand := SetSellTriggerCommand{UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], setSellTriggerCommand)
    
    case "DUMPLOG":
      if(len(parts) < 3) {
        break
      }
      stockSymbol := parts[2]
      dumplogCommand := DumplogCommand{UserId: userId, StockSymbol: stockSymbol}
      userMap[userId] = append(userMap[userId], dumplogCommand)
     
    case "DISPLAY_SUMMARY":
      displaySummary := DisplaySummary{UserId : userId}
      userMap[userId] = append(userMap[userId], displaySummary)
    }
  }
  
  handleCommand(userMap)
}
