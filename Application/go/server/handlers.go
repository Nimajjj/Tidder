package server

import (
  "html/template"
  "net/http"

  SQL "github.com/Nimajjj/Tidder/go/sql"
)

/*
  IndexHandler(db *SQL.SqlServer)

  Function handling the index.html template
*/
func IndexHandler(db *SQL.SqlServer) {
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
}
