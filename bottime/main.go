package main

import (
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "time"
  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sectostr(t int64) string { //преобразуем время в секундах к типу "xx days xx hours xx min xx sec"
  var text string = ""
  a := t / (60*60*24)
  if a != 0 {
    text = fmt.Sprintf("%d days ", a)
  }
  text += fmt.Sprintf("%d hours %d min %d sec", t/(60*60) % 24, t/60 % 60, t%60)
  return text
}

func opendb(dbchan chan *sql.DB) { // открытие базы данных и создание таблиц
  db, err := sql.Open("sqlite3", "./mydatabase.db")
  if err != nil {
    log.Fatal(err)
  }
  log.Println("db Open")

  statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, user_id INTEGER, chat_id INTEGER)")
  statement.Exec()

  statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS times (id INTEGER PRIMARY KEY, time INTEGER, user_id INTEGER, FOREIGN KEY (user_id) REFERENCES users (id))")
  statement.Exec()
  dbchan <- db
}

func updatedb(newdata chan [3]int64, db *sql.DB) { // функция для паралельного обновления данных
  var data [3]int64
  defer db.Close()
  for {
    data = <- newdata
    rows, err := db.Query("SELECT id FROM users WHERE chat_id = (?)", data[1])
    if err != nil {
      log.Println(err)
      continue
    }
    var id int
    if rows.Next() {
      rows.Scan(&id)
    } else {
      statement, err := db.Prepare("INSERT INTO users (user_id, chat_id) VALUES (?, ?)")
      if err != nil {
        log.Println(err)
        continue
      } else {
        statement.Exec(data[0], data[1])
      }
    }
    rows.Close()
    statement, err := db.Prepare("INSERT INTO times (time, user_id) VALUES (?, ?)")
    if err != nil {
      log.Println(err)
    } else {
      statement.Exec(data[2], id)
    }
  }
}

func processinganswer(newdata chan [3]int64, timesusers map[int64][]int64, update tgbotapi.Update, bot *tgbotapi.BotAPI) { // Обработка сообщения пользователя
  log.Printf("[%d] %s", update.Message.Chat.ID, update.Message.Text)

  switch update.Message.Command() {
  case "add":
    t := time.Now().Unix()
    newdata <- [3]int64{update.Message.From.ID, update.Message.Chat.ID, t}
    if timesusers[update.Message.Chat.ID] == nil {
      timesusers[update.Message.Chat.ID] = []int64{t}
    } else {
      timesusers[update.Message.Chat.ID] = append(timesusers[update.Message.Chat.ID], t)
    }
  case "sum":
    if timesusers[update.Message.Chat.ID] != nil {
      var oldt, s, t int64
      var i int
      s = 0
      for i, t = range timesusers[update.Message.Chat.ID] {
        if t % 2 == 0 {
          oldt = t
        } else {
          s += t - oldt
        }
      }
      
      var text string = ""
      if i % 2 == 0 {
        s += time.Now().Unix() - oldt
        text = " ->"
      }
      
      log.Println(s)

      msg := tgbotapi.NewMessage(update.Message.Chat.ID, sectostr(s) + text)
      bot.Send(msg)
    }
  default:
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What??")
    bot.Send(msg)
  }
}

func main() {
  var dbchan chan *sql.DB = make(chan *sql.DB)

  go opendb(dbchan) // открываем базу данных в паралельном потоке

  bot, err := tgbotapi.NewBotAPI("1763199303:AAHm07HCRcXvqeo_BYd_G7MSxeZN74doZFg") // открываем соединение с telegram
  if err != nil {
    log.Panic(err)
  }

  bot.Debug = true

  log.Printf("Authorized on account %s", bot.Self.UserName)
  db := <- dbchan // дожидаемся открытие базы данных

  rows, _ := db.Query("SELECT users.chat_id, times.time FROM users, times WHERE times.user_id = users.id") // считываем данные
  var chatid, time int64
  var timesusers map[int64][]int64 = make(map[int64][]int64)
  for rows.Next() {
    rows.Scan(&chatid, &time)
    if timesusers[chatid] == nil {
      timesusers[chatid] = []int64{time}
    } else {
      timesusers[chatid] = append(timesusers[chatid], time)
    }
  }
  rows.Close()

  var newdata chan [3]int64 = make(chan [3]int64) // канал для отправки данных на добавления в бд
  go updatedb(newdata, db)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates := bot.GetUpdatesChan(u)
  for update := range updates {
    if update.Message != nil {
      go processinganswer(newdata, timesusers, update, bot)
    }
  }
}
