package server

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"mime/multipart"
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
	cookie, _ := r.Cookie("session_id") // get cookie

	if cookie == nil || cookie.Value == "" { // if not connected
		(*viewData).Connected = false
		return -1
	}

	// if connected
	(*viewData).Connected = true
	(*viewData).Account = db.GetAccountFromSession(cookie.Value)

	if (*viewData).Account.Id == 0 {
		Util.Warning("Account ID == 0 found ! We have fuckng problem captain !")
		(*viewData).Connected = false
		return -1
	}

	if (*viewData).Account.ProfilePicture == "" {
		(*viewData).Account.ProfilePicture = SQL.DefaultPP()
	}

	return (*viewData).Account.Id
}

func callTemplate(templateName string, viewdata *SQL.MasterVD, w http.ResponseWriter) error {
	funcMap := template.FuncMap{
		"blobToUrl": func(u string) template.URL {
			return template.URL(u)
		},
	}

	templates := template.New("").Funcs(funcMap)

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

/*
  IndexHandler(db *SQL.SqlServer)

  Function handling the index.html template
*/
func IndexHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Page = "index"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		viewData.IndexVD.Posts = db.GenerateFeed(IAM)

		// Vote
		type FetchQuery struct {
			IdPost int `json:"id_post"`
			Score  int `json:"score"`
		}
		fetchQuery := &FetchQuery{}
		json.NewDecoder(r.Body).Decode(fetchQuery)

		if fetchQuery.IdPost != 0 && fetchQuery.Score != 0 {
			db.Vote(fetchQuery.IdPost, fetchQuery.Score, IAM)
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
		testConnection(r, &viewData, db)

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

func SignupHandler(db *SQL.SqlServer) { // TODO : handle when user is stupid
	viewData := SQL.MasterVD{}
	viewData.Page = "signup"

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if testConnection(r, &viewData, db) != -1 { // if you're connected how tf did you get here
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
		if testConnection(r, &viewData, db) != -1 { // if you're connected how tf did you get here
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

	http.HandleFunc("/new/post", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		if IAM == -1 { // if you're not connected redirect to home page
			http.Redirect(w, r, "/", http.StatusFound)
		}

		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

		title := r.FormValue("title")
		content := r.FormValue("content")
		nsfw := false
		id_subject := r.FormValue("subtidder")
		id_author := IAM

		// load image
		if r.Method == "POST" {
			var media string
			file, header, err := r.FormFile("media_file")
			if header != nil { // if header != nil then there is a file
				defer func(file multipart.File) {
					err := file.Close()
					if err != nil {
						Util.Error(err)
					}
				}(file)
			}
			if err != nil {
				if err != http.ErrMissingFile {
					Util.Error(err)
				}
			} else {
				media = "data:" + header.Header.Get("Content-Type") + ";base64," // file to base64

				bytes, err := io.ReadAll(file)
				if err != nil {
					Util.Error(err)
				}

				media += base64.StdEncoding.EncodeToString(bytes)
			}

			if title != "" && (content != "" || media != "") {
				db.CreatePost(title, media, content, nsfw, id_subject, id_author)
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

func CreateSubtidderHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Page = "create_subtidder"

	http.HandleFunc("/new/t", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		if IAM == -1 { // if you're not connected redirect to home page
			http.Redirect(w, r, "/", http.StatusFound)
		}

		if r.Method == "POST" {
			subTidderName := r.FormValue("name")
			if subTidderName == "" {
				viewData.Errors.CreateSubtidder = "Please enter a valid subtidder name..."
			} else {
				subTidderNsfw := false
				viewData.Errors.CreateSubtidder = db.CreateSub(subTidderName, IAM, subTidderNsfw)
				if viewData.Errors.CreateSubtidder == "" {
					http.Redirect(w, r, "/t/"+subTidderName, http.StatusFound)
				}
			}
		}

		err := callTemplate("create_subtidder", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func DisconnectHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	http.HandleFunc("/disconnect", func(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: "", Expires: expiration, Path: "/"}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)

		err := callTemplate("disconnect", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func PostHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	viewData.Page = "post_page"

	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		var postVD SQL.PostVD
		var ID = postVD.Post.Id

		post := db.GetPosts("id_post=" + strconv.Itoa(ID))[0]
		postVD.Post = db.MakeDisplayablePost(post, IAM)

		err := callTemplate("post_page", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func ProfilePageHandler(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	http.HandleFunc("/u/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		viewData.CreatePostsVD.SubscribedSubjects = db.GetSubtiddersSubscribed(IAM)

		err := callTemplate("profile_page", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}
