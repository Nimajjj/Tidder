package utility

import (
    "fmt"
    "os"
    "time"
    "log"
)


func Log(text string) {
  logs := time.Now().Format("01-02-2006 15:04:05") + " \t" + text
  addToLogs(logs, false)
  fmt.Println(logs)
}


func Error(err error) {
  addToLogs(err.Error(), true)
  log.Fatal(err)
}


func addToLogs(text string, isErr bool) {
  file, err := os.OpenFile("./logs/master.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil { panic(err) }
  defer file.Close()

  if isErr {
    _, err2 := file.WriteString("\n" + time.Now().Format("01-02-2006 15:04:05") + " \t" + text)
    if err2 != nil { panic(err2) }
  } else {
    _, err2 := file.WriteString("\n" + text)
    if err2 != nil { panic(err2) }
  }


}
