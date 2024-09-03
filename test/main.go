package main

import (
  "fmt"
  "github.com/ssimunic/gosensors"
  "strings"
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
			if key == "temp1" {
        var temp float64
        fmt.Sscan(strings.ReplaceAll(strings.ReplaceAll(value, "°C", ""), "+", ""), &temp)
  			fmt.Println(key, temp)
      }
		}
	}
}
