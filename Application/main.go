package main

import (
  Server "github.com/Nimajjj/Tidder/go/server"
  "os"
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

func main() {
	DatabaseIp := "127.0.0.1"
	if len(os.Args) > 1 {
		DatabaseIp = os.Args[1]
	}
	Server.Run(DatabaseIp)
}
