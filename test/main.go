package main

import (
  "fmt"
//  "database/sql"
//  _ "github.com/mattn/go-sqlite3"
//  "log"
//  "time"
//  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func F(c chan int) {
  var a int
  fmt.Scan(&a)
  c <- a
}

func main() {
  var c chan int = make(chan int)
  go F(c)
  fmt.Println("Ждём ответ от другого потока")
  fmt.Println("Получили ответ", <-c)
}
