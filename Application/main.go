package main

import (
  Main "github.com/Nimajjj/Tidder/go"
  Global "github.com/Nimajjj/Tidder/go/global"

  "fmt"
)

func main()  {
  fmt.Println("\nTidder Inc © 2022. Tous droits réservés")
  fmt.Println("Starting server : http://localhost:80\n")

  Main.TestFunction()
  Global.TestFunction()

  fmt.Println()
}
