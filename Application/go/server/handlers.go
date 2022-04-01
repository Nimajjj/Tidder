package server

import (
  "encoding/json"
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
  viewData := SQL.MasterVD{}

  tpl := template.Must(template.ParseFiles("./pages/index.html"))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.FormValue("name") != "" {        // create subtidder
      subTidderName := r.FormValue("name")
      subTidderNsfw := false
      if r.FormValue("nsfw") == "0" {
        subTidderNsfw = true
      }
      viewData.Error = db.CreateSub(subTidderName, 2, subTidderNsfw)
    }

    tpl.Execute(w, viewData)
  })
}


func SubtidderHandler(db *SQL.SqlServer) {
  var subtidder SQL.SubtidderViewData

  tpl := template.Must(template.ParseFiles("./pages/subtidder.html"))
  http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
    // CREATE SUBTIDDER COMPONENT //////////////////////////////////////////////////

    if r.FormValue("name") != "" {       
      subTidderName := r.FormValue("name")
      subTidderNsfw := false
      if r.FormValue("nsfw") == "0" {
        subTidderNsfw = true
      }
      db.CreateSub(subTidderName, 1, subTidderNsfw)
    }

    // END CREATE SUBTIDDER COMPONENT //////////////////////////////////////////////
    // SUBSCRIBE TO SUBTIDDER COMPONENT ////////////////////////////////////////////

    type Subscription struct {
      IdAccount int `json:"id_account_subscribing"`
      IdSubject int `json:"id_subject_to_subscribe"`
    }
    subscription := &Subscription{}
    json.NewDecoder(r.Body).Decode(subscription)
    if subscription.IdAccount != 0 && subscription.IdSubject != 0 {
      db.SubscribeToSubject(subscription.IdAccount, subscription.IdSubject)
    }

    // END SUBSCRIBE TO SUBTIDDER COMPONENT ////////////////////////////////////////
    // MAIN SUBTIDDER COMPONENT ////////////////////////////////////////////////////
    
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
    // END MAIN SUBTIDDER COMPONENT ////////////////////////////////////////////////
  })
}


func SearchHandler(db *SQL.SqlServer) {
  results := SQL.SearchViewData{}
  tpl := template.Must(template.ParseFiles("./pages/search/search.html"))
  http.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
    results.Subjects = map[SQL.Subject]int{}
    search := ""
    if r.FormValue("search") != "" { search = r.FormValue("search") }

    for _, subject := range db.GetSubs("name LIKE \"%" + search + "%\"") {
      results.Subjects[subject] = db.GetNumberOfSubscriber(subject.Id)
    }
    tpl.Execute(w, results)
  })
}