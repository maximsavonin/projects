package main

import (
  "fmt"
//  "database/sql"
//  _ "github.com/mattn/go-sqlite3"
//  "log"
//  "time"
//  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func F(a map[int][]int){
  a[1] = []int{10}
}

func main() {
  var x map[int][]int = make(map[int][]int)
  F(x)
  fmt.Println(x[1])
}
