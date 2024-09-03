package main

import (
  "fmt"
  "github.com/ssimunic/gosensors"
//  "database/sql"
//  _ "github.com/mattn/go-sqlite3"
//  "log"
//  "time"
//  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
  sensors, err := gosensors.NewFromSystem()
  if err != nil {
    panic(err)
  }

  fmt.Println(sensors)

  for chip := range sensors.Chips {
		// Iterate over entries
		for key, value := range sensors.Chips[chip] {
			// If CPU or GPU, print out
			fmt.Println(key, value)
		}
	}
}
