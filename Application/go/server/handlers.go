package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	SQL "github.com/Nimajjj/Tidder/go/sql"
	Util "github.com/Nimajjj/Tidder/go/utility"
)

/*
	return -1 == not connected
	return 0 == connected
*/
func testConnection(r *http.Request, viewData *SQL.MasterVD, db *SQL.SqlServer) int {
	cookie, _ := r.Cookie("session_id")
	if cookie == nil {
		Util.Log("No cookie found")
		return -1
	}
	(*viewData).Connected = true
	(*viewData).Account = db.GetAccountFromSession(cookie.Value)
	return (*viewData).Account.Id
}

func callTemplate(templateName string, viewdata *SQL.MasterVD, w http.ResponseWriter) error {
	templates := template.New("")
	templates, err := templates.ParseFiles("./pages/base.html", "./pages/templates/"+templateName+".html")
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob("./pages/templates/components/*.html")
	if err != nil {
		return err
	}

	err = templates.ExecuteTemplate(w, "base", *viewdata)
	if err != nil {
		return err
	}

	_ = templateName
	return nil
}

func popup(w http.ResponseWriter, r *http.Request, viewData *SQL.MasterVD, db *SQL.SqlServer, IAM int) {
	if r.FormValue("name") != "" { // SUBTIDDER CREATION //
		subTidderName := r.FormValue("name")
		subTidderNsfw := false
		if r.FormValue("nsfw") == "0" {
			subTidderNsfw = true
		}
		(*viewData).Errors.CreateSubtidder = db.CreateSub(subTidderName, 2, subTidderNsfw)
	} else if r.FormValue("SubmitSignin") != "" { // SIGN IN //
		username := r.FormValue("input_pseudo_signin")
		password := r.FormValue("input_password_signin")
		connectedUsr, sessionId := db.TryToConnectUser(username, password)
		if sessionId != "" {
			(*viewData).Connected = true
			(*viewData).Account = connectedUsr[0]
			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{Name: "session_id", Value: sessionId, Expires: expiration, Path: "/"}
			http.SetCookie(w, &cookie)
		}
	} else if r.FormValue("submit_post") == "Envoyer" { // POST CREATION //
		title := r.FormValue("post_title")
		media_url := ""
		content := r.FormValue("post_content")
		nsfw := false
		id_subject := r.FormValue("post_subtidder")
		id_author := IAM

		if content != "" && title != "" {
			db.CreatePost(title, media_url, content, nsfw, id_subject, id_author)
		} else {
			Util.Warning("An error occured during post creation.")
			(*viewData).Errors.CreatePost = "An error occured during post creation."
		}

	} else { // SIGN UP //
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

/*
  IndexHandler(db *SQL.SqlServer)

  Function handling the index.html template
*/
func IndexHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		popup(w, r, &viewData, db, IAM)
		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)
		viewData.IndexVD.Posts = db.GenerateFeed(IAM)

		err := callTemplate("index_feed", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func SubtidderHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	var subtidder SQL.SubtidderViewData

	http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		popup(w, r, &viewData, db, IAM)
		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

		// MAIN SUBTIDDER COMPONENT //
		id := strings.ReplaceAll(r.URL.Path, "localhost/t/", "")
		id = strings.ReplaceAll(r.URL.Path, "/t/", "")

		subtidder.Posts = []map[SQL.Posts]SQL.Accounts{}
		subtidder.Sub = db.GetSubs("name=\"" + id + "\"")[0] // ca c'est sale -> a refaire

		posts := db.GetPosts("id_subject=" + strconv.Itoa(subtidder.Sub.Id) + " ORDER BY creation_date DESC")

		for _, post := range posts { // svp me demandez pas d'expliquer ce bout de code franchement c chaud
			subtidder.Posts = append(
				subtidder.Posts,
				map[SQL.Posts]SQL.Accounts{
					post: db.GetAccountById(post.IdAuthor),
				},
			)
		}

		// SUBSCRIBE TO SUBTIDDER COMPONENT //
		type Subscription struct {
			IdAccount int `json:"id_account_subscribing"`
			IdSubject int `json:"id_subject_to_subscribe"`
		}
		subscription := &Subscription{}
		json.NewDecoder(r.Body).Decode(subscription)
		if subscription.IdAccount != 0 && subscription.IdSubject != 0 {
			db.SubscribeToSubject(IAM, subtidder.Sub.Id)
		}

		subtidder.Subscribed = db.IsSubscribeTo(IAM, subtidder.Sub.Id)
		viewData.SubtidderVD = subtidder

		err := callTemplate("subtidder", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func SearchHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}

	http.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		popup(w, r, &viewData, db, IAM)
		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

		// SEARCH COMPONENT //
		viewData.SearchVD.Subjects = map[SQL.Subject]int{}
		search := ""
		if r.FormValue("search") != "" {
			search = r.FormValue("search")
		}

		for _, subject := range db.GetSubs("name LIKE \"%" + search + "%\"") {
			viewData.SearchVD.Subjects[subject] = db.GetNumberOfSubscriber(subject.Id)
		}

		err := callTemplate("search", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}
