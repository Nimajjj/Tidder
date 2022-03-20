package main

import (
  Main "github.com/Nimajjj/Tidder/go"
  Server "github.com/Nimajjj/Tidder/go/server"

  "fmt"
)

/*
  Must stay as clean as possible ! (jvous regarde vous deux et vos codes jamais formatés (-_-) )
*/

func main()  {
  checkImport()
  Server.Run()
}

/*
  Exemple func, to remove
*/
func checkImport()  {
  Main.TestFunction()
  fmt.Println()
}
