package utility

import (
	"fmt"
	"log"
	"os"
	"time"
)

/*
  Log(text string)

  Format log (dd-mm-yyyy hh-mm-ss   <textual log>)
  Print log in terminal
  Print log in ./logs/master.log

*/
func Log(text string) {
	logs := time.Now().Format("01-02-2006 15:04:05") + " \tDEBUG\t\t" + text
	registerLog(logs)
	fmt.Println(logs)
}

/*
  Error(err error)

  Same as Log(text string) but handling error
  Kill program
*/
func Error(err error) {
	logs := time.Now().Format("01-02-2006 15:04:05") + " \tERROR /!\\\t\t" + err.Error()
	registerLog(logs)
	log.Fatal(err)
}

func Warning(text string) {
	logs := time.Now().Format("01-02-2006 15:04:05") + " \tWARNING\t\t" + text
	registerLog(logs)
	fmt.Println(logs)
}

func Query(querry string) {
	logs := time.Now().Format("01-02-2006 15:04:05") + " \tQUERRY\t\t" + querry
	registerLog(logs)
	fmt.Println(logs)
}

/*
  addToLogs(text string, isErr bool)

  Print text in ./logs/master.log
*/
func registerLog(text string) {
	file, err := os.OpenFile("./logs/master.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err2 := file.WriteString("\n" + text)
	if err2 != nil {
		panic(err2)
	}
}
