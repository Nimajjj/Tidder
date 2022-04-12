package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	SQL "github.com/Nimajjj/Tidder/go/sql"
)

/*
  IndexHandler(db *SQL.SqlServer)

  Function handling the index.html template
*/
func IndexHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Connected = true
	viewData.CreatePostsVD = SQL.CreatePostsVD{}

	IAM := 1

	viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

	tpl := template.Must(template.ParseFiles("./pages/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// CREATE SUBTIDDER COMPONENT //
		if r.FormValue("name") != "" {
			subTidderName := r.FormValue("name")
			subTidderNsfw := false
			if r.FormValue("nsfw") == "0" {
				subTidderNsfw = true
			}
			viewData.Errors.CreateSubtidder = db.CreateSub(subTidderName, 2, subTidderNsfw)
		}

		// SIGN UP COMPONENT //
		if r.Method == http.MethodPost {
			
			if (r.FormValue("submit_post") == "Envoyer") {
				title := r.FormValue("post_title")
				media_url := ""
				content := r.FormValue("post_content")
				nsfw := false
				id_subject := r.FormValue("post_subtidder")
				id_author := IAM

				db.CreatePost(title, media_url, content, nsfw, id_subject, id_author)
			} else {
				pseudo := r.FormValue("pseudo_input")
				email := r.FormValue("email_input")
				password := r.FormValue("password_input")
				verifpassword := r.FormValue("passwordverif_input")
				birthdate := r.FormValue("birthdate_input")
				studentId := r.FormValue("id_input")
				if pseudo != "" && email != "" && password != "" && birthdate != "" && studentId != "" {
					viewData.Errors.Signup = db.CreateAccount(pseudo, email, password, birthdate, studentId, verifpassword)
				}
			}
		}

		tpl.Execute(w, viewData)
	})
}

func SubtidderHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Connected = true
	var subtidder SQL.SubtidderViewData
	

	tpl := template.Must(template.ParseFiles("./pages/subtidder.html"))
	http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
		// CREATE SUBTIDDER COMPONENT //
		if r.FormValue("name") != "" {
			subTidderName := r.FormValue("name")
			subTidderNsfw := false
			if r.FormValue("nsfw") == "0" {
				subTidderNsfw = true
			}
			db.CreateSub(subTidderName, 1, subTidderNsfw)
		}

		// SUBSCRIBE TO SUBTIDDER COMPONENT //

		type Subscription struct {
			IdAccount int `json:"id_account_subscribing"`
			IdSubject int `json:"id_subject_to_subscribe"`
		}
		subscription := &Subscription{}
		json.NewDecoder(r.Body).Decode(subscription)
		if subscription.IdAccount != 0 && subscription.IdSubject != 0 {
			db.SubscribeToSubject(subscription.IdAccount, subscription.IdSubject)
		}

		// MAIN SUBTIDDER COMPONENT //

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
		subtidder.Sub = db.GetSubs("name=\"" + formated_id + "\"")[0] // ca c'est sale -> a refaire

		posts := db.GetPosts("id_subject=" + strconv.Itoa(subtidder.Sub.Id) + " ORDER BY creation_date DESC")

		for _, post := range posts { // svp me demandez pas d'expliquer ce bout de code franchement c chaud
			subtidder.Posts = append(
				subtidder.Posts,
				map[SQL.Posts]SQL.Accounts{
					post: db.GetAccountById(post.IdAuthor),
				},
			)
		}

		viewData.SubtidderVD = subtidder
		tpl.Execute(w, viewData)
	})
}

func SearchHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Connected = true

	tpl := template.Must(template.ParseFiles("./pages/search/search.html"))
	http.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		// SEARCH COMPONENT //
		viewData.SearchVD.Subjects = map[SQL.Subject]int{}
		search := ""
		if r.FormValue("search") != "" {
			search = r.FormValue("search")
		}

		for _, subject := range db.GetSubs("name LIKE \"%" + search + "%\"") {
			viewData.SearchVD.Subjects[subject] = db.GetNumberOfSubscriber(subject.Id)
		}

		tpl.Execute(w, viewData)
	})
}

func ConnectionHandler(db *SQL.SqlServer) {

	Testc := SQL.MasterVD{}
	tpltest := template.Must(template.ParseFiles("./test/connectiontest.html"))
	http.HandleFunc("/testco/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			pseudo := r.FormValue("pseudo_input")
			password := r.FormValue("password_input")

			db.Connection(Testc.Account.Id, pseudo, password, Testc.Connected)
		}
		tpltest.Execute(w, Testc)
	})
}
