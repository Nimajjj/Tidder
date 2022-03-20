package main

import (
  Main "github.com/Nimajjj/Tidder/go"
  Server "github.com/Nimajjj/Tidder/go/server"

  "fmt"
)

/*
  Tidder Inc © 2022. Tous droits réservés

  Ynov Aix - 2021/2022
  Created : 20/03/2022
  Authors :
    JEHAM   Laurie
    OBRY    Maxime
    BORELLO Benjamin
*/

/*
  /!\ MUST STAY AS CLEAN AS POSSIBLE /!\ (jvous vois vous deux et vos codes jamais formatés (-_-) )
  /!\ All func must be commented (cf. "/go/server/server.go") /!\
*/

func main()  {
  checkImport()
  Server.Run()
}

/*
  Exemple func, to remove :O
*/
func checkImport()  {
  Main.TestFunction()
  fmt.Println()
}
