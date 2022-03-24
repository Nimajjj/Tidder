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

  tpl := template.Must(template.ParseFiles("./pages/index.html"))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.FormValue("name") != "" {
      subTidderName := r.FormValue("name")
      subTidderNsfw := false
      if r.FormValue("nsfw") == "1" {
        subTidderNsfw = true
      }
      db.CreateSub(subTidderName, 2, subTidderNsfw)
    }

    tpl.Execute(w, account)
  })
}


func SubtidderHandler(db *SQL.SqlServer) {
  var subtidder SQL.Subtidder

  tpl := template.Must(template.ParseFiles("./pages/subtidder.html"))
  http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
    id := strings.ReplaceAll(r.URL.Path, "localhost/t/", "")
		id = strings.ReplaceAll(r.URL.Path, "/t/", "")
    subtidder.Sub = db.GetSubs("name="+id)[0]  // ca c'est sale -> a refaire

    tpl.Execute(w, subtidder)
  })
}
