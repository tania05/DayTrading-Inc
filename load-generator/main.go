package main

import "fmt"
import "io/ioutil"
import "strings"

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

func (command *BuyCommand) name() string {
  return fmt.Sprintf("userid %s stocksymbol %s amount %d", command.userId, command.stockSymbol, command.amount)
}

func main() {
  b, err := ioutil.ReadFile("10_user_workload")

  if (err != nil) {
    panic(err)
  }

  lines := strings.Split(string(b), "\n")

  userMap := make(map[string][]string)
  x := 0
  for _, line := range lines {
    x += 2
    parts := strings.Split(line, ",")
    if len(parts) < 2 {
      continue
    }

    switch commandType := strings.Split(parts[0], ' ')[1]; commandType {

    fmt.Println(parts)
    userId := strings.Trim(parts[1], " ")
  }

  //fmt.Println(userMap)
}
