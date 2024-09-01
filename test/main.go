package main

import (
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "time"
  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sectostr(t int) string {
  var text string = ""
  a := t / (60*60*24)
  if a != 0 {
    text = fmt.Sprintf("%d days", a)
  }
  text += fmt.Sprintf(" %d hours %d min %d sec", t/(60*60) % 24, t/60 % 60, t%60)
  return text
}

func Opendb(dbchan chan *sql.DB) {
  db, err := sql.Open("sqlite3", "./mydatabase.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  log.Println("db Open")

  statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, user_id INTEGER, chat_id INTEGER)")
  statement.Exec()

  statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS times (id INTEGER PRIMARY KEY, time INTEGER, user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users (id))")
  statement.Exec()
  for {
    dbchan <- db
    _ = <- dbchan
  }
}

func F(dbchan chan *sql.DB, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
  log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)
  
  db := <- dbchan
  rows, err := db.Query("SELECT id FROM users WHERE chat_id = (?)", update.Message.Chat.ID)
  if err != nil { 
    log.Println(err) 
  } 
  var id int 
  if rows.Next() { 
    rows.Scan(&id) 
  } else { 
    statement, err := db.Prepare("INSERT INTO users (user_id, chat_id) VALUES (?, ?)") 
    if err != nil {
      log.Println(err)
    }
    statement.Exec(update.Message.Chat.ID, update.Message.From.ID)
  }
  rows.Close()
  switch update.Message.Command() {
  case "add":
    statement, err := db.Prepare("INSERT INTO times (time, user_id) VALUES (?, ?)")
    if err != nil {
      log.Println(err)
    }
    statement.Exec(time.Now().Unix(), id)
    dbchan <- db
  case "sum":
    rowstimes, err := db.Query("SELECT time FROM times WHERE user_id = (?)", id)
    if err != nil {
      log.Println(err)
    }
    defer rowstimes.Close()

    var tfir int
    var tsec int
    var s int = 0
    var i int = 0

    for rowstimes.Next() {
      if i % 2 == 0 {
        rowstimes.Scan(&tfir)
      } else {
        rowstimes.Scan(&tsec)
        s += tsec - tfir
      }
      i += 1
      log.Println(s)
    }
    rowstimes.Close()
    dbchan <- db

    var text string = ""

    if i % 2 == 1 {
      s += int(time.Now().Unix()) - tfir
      text = "  ->"
    }
    
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, sectostr(s) + text)
    bot.Send(msg)
  default:
    dbchan <- db
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What??")
    bot.Send(msg)
  }
}

func main() {
  var db chan *sql.DB = make(chan *sql.DB)

  go Opendb(db)

  bot, err := tgbotapi.NewBotAPI("1763199303:AAHm07HCRcXvqeo_BYd_G7MSxeZN74doZFg")
  if err != nil {
    log.Panic(err)
  }

  bot.Debug = true

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates := bot.GetUpdatesChan(u)
  for update := range updates {
    if update.Message != nil {
      go F(db, update, bot)
    }
  }
}
