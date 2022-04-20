package server

import (
	_ "encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	SQL "github.com/Nimajjj/Tidder/go/sql"
)

func popup(w http.ResponseWriter, r *http.Request, viewData *SQL.MasterVD, db *SQL.SqlServer) {
	IAM := 1
	// CREATE SUBTIDDER COMPONENT //
	if r.FormValue("name") != "" {
		subTidderName := r.FormValue("name")
		subTidderNsfw := false
		if r.FormValue("nsfw") == "0" {
			subTidderNsfw = true
		}
		(*viewData).Errors.CreateSubtidder = db.CreateSub(subTidderName, 2, subTidderNsfw)
	} else if r.FormValue("SubmitSignin") != "" {		// SIGN IN COMPONENT //
		username := r.FormValue("input_pseudo_signin")
		password := r.FormValue("input_password_signin")
		connectedUsr, sessionId := db.TryToConnectUser(username, password)
		if sessionId != "" {
			(*viewData).Connected = true
			(*viewData).Account = connectedUsr[0]
			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{Name: "session_id", Value: sessionId, Expires: expiration}
			http.SetCookie(w, &cookie)
		}
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
				(*viewData).Errors.Signup = db.CreateAccount(pseudo, email, password, birthdate, studentId, verifpassword)
			}
		}
	}
}


/*
  IndexHandler(db *SQL.SqlServer)

  Function handling the index.html template
*/
func IndexHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.CreatePostsVD = SQL.CreatePostsVD{}

	IAM := 1

	viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)


	tpl := template.Must(template.ParseFiles("./pages/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		popup(w, r, &viewData, db)

		tpl.Execute(w, viewData)
	})
}

func SubtidderHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	var subtidder SQL.SubtidderViewData
	

	tpl := template.Must(template.ParseFiles("./pages/subtidder.html"))
	http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
		popup(w, r, &viewData, db)

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

	tpl := template.Must(template.ParseFiles("./pages/search/search.html"))
	http.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		popup(w, r, &viewData, db)
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
