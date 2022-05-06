package server

import (
	"encoding/json"
	"html/template"
	"net/http"
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
		(*viewData).Errors.CreateSubtidder = db.CreateSub(subTidderName, IAM, subTidderNsfw)
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

		// load image

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
	viewData.Page = "index"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		popup(w, r, &viewData, db, IAM)
		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)
		viewData.IndexVD.Posts = db.GenerateFeed(IAM)

		// Vote
		type Vote struct {
			IdPost int `json:"id_post"`
			Score  int `json:"score"`
		}
		vote := &Vote{}
		json.NewDecoder(r.Body).Decode(vote)

		if vote.IdPost != 0 && vote.Score != 0 {
			db.Vote(vote.IdPost, vote.Score, IAM)
		}

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
	viewData.Page = "subtidder"

	http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		popup(w, r, &viewData, db, IAM)
		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

		// MAIN SUBTIDDER COMPONENT //
		id := strings.ReplaceAll(r.URL.Path, "localhost/t/", "")
		id = strings.ReplaceAll(r.URL.Path, "/t/", "")

		subtidder.Sub = db.GetSubs("name=\"" + id + "\"")[0] // ca c'est sale -> a refaire
		subtidder.Posts = db.GenerateSubTidderFeed(IAM, subtidder.Sub.Id)

		type FetchQuery struct {
			IdPost    int `json:"id_post"`
			Score     int `json:"score"`
			IdAccount int `json:"id_account_subscribing"`
			IdSubject int `json:"id_subject_to_subscribe"`
		}
		fetchQuery := &FetchQuery{}
		json.NewDecoder(r.Body).Decode(fetchQuery)

		// VOTES //
		if fetchQuery.IdPost != 0 && fetchQuery.Score != 0 && IAM != -1 {
			db.Vote(fetchQuery.IdPost, fetchQuery.Score, IAM)
		}

		// SUBSCRIBE TO SUBTIDDER COMPONENT //
		if fetchQuery.IdAccount != 0 && fetchQuery.IdSubject != 0 {
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
	viewData.Page = "search"

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

func SignupHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Page = "signup"

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if IAM != -1 { // if you're connected how tf did you get here
			http.Redirect(w, r, "/", http.StatusFound)
		}
		pseudo := r.FormValue("pseudo_input")
		email := r.FormValue("email_input")
		password := r.FormValue("password_input")
		verifpassword := r.FormValue("passwordverif_input")
		birthdate := r.FormValue("birthdate_input")
		if r.Method == "POST" {
			viewData.Errors.Signup = db.CreateAccount(pseudo, email, password, birthdate, "", verifpassword)
			if viewData.Errors.Signup == "" {
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}

		err := callTemplate("signup", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func SigninHandler(db *SQL.SqlServer) { // TODO : handle when user is stupid
	viewData := SQL.MasterVD{}
	viewData.Page = "signup"

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		if IAM != -1 { // if you're connected how tf did you get here
			http.Redirect(w, r, "/", http.StatusFound)
		}
		if r.Method == "POST" {
			username := r.FormValue("pseudo_input")
			password := r.FormValue("password_input")
			connectedUsr, sessionId := db.TryToConnectUser(username, password)

			if sessionId != "" {
				viewData.Connected = true
				viewData.Account = connectedUsr[0]
				expiration := time.Now().Add(24 * time.Hour)
				cookie := http.Cookie{Name: "session_id", Value: sessionId, Expires: expiration, Path: "/"}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}

		err := callTemplate("signin", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func CreatePostHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Page = "create_post"

	http.HandleFunc("/create_post", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)

		if IAM == -1 { // if you're not connected redirect to home page
			http.Redirect(w, r, "/signin", http.StatusFound)
		}

		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

		title := r.FormValue("title")
		media_url := ""
		content := r.FormValue("content")
		nsfw := false
		id_subject := r.FormValue("subtidder")
		id_author := IAM

		// load image
		if r.Method == "POST" {
			if content != "" && title != "" {
				db.CreatePost(title, media_url, content, nsfw, id_subject, id_author)
				http.Redirect(w, r, "/", http.StatusFound) // TODO : redirect to post page
			} else {
				Util.Warning("An error occured during post creation.")
				viewData.Errors.CreatePost = "An error occured during post creation."
			}
		}

		err := callTemplate("create_post", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}
