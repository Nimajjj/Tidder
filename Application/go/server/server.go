package server

import (
  "fmt"
  "html/template"
	"net/http"

  SQL "github.com/Nimajjj/Tidder/go/sql"
)

/*
  /go/server/server.go

  All files path must be relative to the main.go script position.
  Only functions starting with a CAPITAL LETTER can be access from the main.go script
*/


/*
  Run()

  ONLY function which can be accessed from the main.go script
  Entry door for the go web server.
*/
func Run()  {
  fmt.Println("\nTidder Inc © 2022. Tous droits réservés")
  fmt.Println("Starting server : http://localhost:80")

  initStaticFolders()
  launchServer()
}

/*
  initStaticFolders()

  All statics folders which will be used in html/css/js files must be declared here.
*/
func initStaticFolders() {
  cssFolder := http.FileServer(http.Dir("./style"))
  imgFolder := http.FileServer(http.Dir("./images"))
  jsFolder := http.FileServer(http.Dir("./scripts"))
  http.Handle("/style/", http.StripPrefix("/style/", cssFolder))
  http.Handle("/images/", http.StripPrefix("/images/", imgFolder))
  http.Handle("/scripts/", http.StripPrefix("/scripts/", jsFolder))
}

/*
  launchServer()

  Main function which run the server
  For the moment the template & the http.ListenAndServe(":80", nil) are here.

  To do :
    -create an individual function for each template.
*/
func launchServer() {
  var db SQL.SqlServer
  db.Connect()
  defer db.Close()
  account := db.GetAccountById(1)

  indexTpl := template.Must(template.ParseFiles("./pages/index.html"))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.FormValue("name") != "" {
      subTidderName := r.FormValue("name")
      subTidderNsfw := false
      if r.FormValue("nsfw") == "1" {
        subTidderNsfw = true
      }
      db.CreateSub(subTidderName, 2, subTidderNsfw)
    }

		indexTpl.Execute(w, account)
	})
  http.ListenAndServe(":80", nil)
}
