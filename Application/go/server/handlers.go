package server

import (
  "strconv"
  "strings"
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
    if r.FormValue("name") != "" {        // create subtidder
      subTidderName := r.FormValue("name")
      subTidderNsfw := false
      if r.FormValue("nsfw") == "0" {
        subTidderNsfw = true
      }
      db.CreateSub(subTidderName, 2, subTidderNsfw)
    }

    tpl.Execute(w, account)
  })
}


func SubtidderHandler(db *SQL.SqlServer) {
  var subtidder SQL.SubtidderViewData

  tpl := template.Must(template.ParseFiles("./pages/subtidder.html"))
  http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
    if r.FormValue("name") != "" {        // create subtidder
      subTidderName := r.FormValue("name")
      subTidderNsfw := false
      if r.FormValue("nsfw") == "0" {
        subTidderNsfw = true
      }
      db.CreateSub(subTidderName, 2, subTidderNsfw)
    }

    
    id := strings.ReplaceAll(r.URL.Path, "localhost/t/", "")
		id = strings.ReplaceAll(r.URL.Path, "/t/", "")

    formated_id := ""
    for _, char := range id {
      if char == '+' {
        formated_id += " "
      } else {
        formated_id += string(char)
      }
    }

    subtidder.Posts = []map[SQL.Posts]SQL.Accounts{}
    subtidder.Sub = db.GetSubs("name=\"" + formated_id + "\"")[0]  // ca c'est sale -> a refaire

    posts := db.GetPosts("id_subject=" + strconv.Itoa(subtidder.Sub.Id) + " ORDER BY creation_date DESC")

    for _, post := range posts {  // svp me demandez pas d'expliquer ce bout de code franchement c chaud
      subtidder.Posts = append(
        subtidder.Posts,
        map[SQL.Posts]SQL.Accounts{
          post: db.GetAccountById(post.IdAuthor),
        },
      )
    }

    tpl.Execute(w, subtidder)
  })
}
