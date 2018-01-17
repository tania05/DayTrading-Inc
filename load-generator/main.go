package main

import "fmt"
import "io/ioutil"
import "strings"
import "strconv"

//type AddCommand struct {
//  userId string
//  amount int
//}
//
//type QuoteCommand struct {
//  userId string
//  stockSymbol string
//}

type BuyCommand struct {
  userId string
  stockSymbol string
  amount int
}

type Command interface {
  name() string
}

func (command BuyCommand) name() string {
  return fmt.Sprintf("userid %s stocksymbol %s amount %d", command.userId, command.stockSymbol, command.amount)
}

func main() {
  b, err := ioutil.ReadFile("10_user_workload")

  if (err != nil) {
    panic(err)
  }

  lines := strings.Split(string(b), "\n")
  //// "kk", []
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
    
    case "BUY":
      a64, _ := strconv.ParseInt(strings.Replace(parts[3],".","",-1),10,0)
      var amount = int(a64)
      stockSymbol := parts[2]
      
      buyCommand := BuyCommand{userId: userId, stockSymbol: stockSymbol, amount: amount}
      userMap[userId] = append(userMap[userId], buyCommand)
      // fmt.Println(buyCommand)
    }
  }

  // fmt.Println(userMap)
}
