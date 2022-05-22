package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func getSuffix(r *http.Request) string {
	if strings.Contains(r.URL.Path, "/") {
		split := strings.Split(r.URL.Path, "/")
		if split[len(split)-1] != "" {
			return split[len(split)-1]
		} else {
			return "!err404"
		}
	}
	return "!err404"
}

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

	if (*viewData).Account.ProfilePicture == "default" || (*viewData).Account.ProfilePicture == "" {
		(*viewData).Account.ProfilePicture = SQL.DefaultPP()
	}

	return (*viewData).Account.Id
}

func callTemplate(templateName string, viewdata *SQL.MasterVD, w http.ResponseWriter) error {
	funcMap := template.FuncMap{
		"blobToUrl": func(u string) template.URL {
			return template.URL(u)
		},
		"hasAccess": func(access SQL.SubjectAccess, accessList string) bool {
			accessSplit := strings.Split(accessList, ",")
			for _, a := range accessSplit {
				fmt.Println(a)
				switch a {
				case "pin_post":
					if access.Pin == 1 {
						return true
					}
				case "create_post":
					if access.CreatePost == 1 {
						return true
					}
				case "manage_post":
					if access.ManagePost == 1 {
						return true
					}
				case "manage_subject":
					if access.ManageSub == 1 {
						return true
					}
				case "manage_role":
					if access.ManageRole == 1 {
						return true
					}
				case "ban_user":
					if access.BanUser == 1 {
						return true
					}
				}
			}
			return false
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

		id := getSuffix(r)
		if id == "!err404" {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}

		// MAIN SUBTIDDER COMPONENT //
		subs := db.GetSubs("name=\"" + id + "\"")
		if len(subs) == 0 {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}
		subtidder.Sub = subs[0]
		subtidder.Roles = db.GenerateRoleAccess(subtidder.Sub.Id)
		subtidder.UserRole = db.GenerateUserRoleAccess(subtidder.Sub.Id, IAM)

		if r.Method == "POST" {
			if true {
				var pp string
				file, header, err := r.FormFile("pp_input")
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
					pp = "data:" + header.Header.Get("Content-Type") + ";base64," // file to base64

					bytes, err := io.ReadAll(file)
					if err != nil {
						Util.Error(err)
					}

					pp += base64.StdEncoding.EncodeToString(bytes)
				}

				if pp != "" {
					db.UpdateSubtidder("profile_picture", pp, subtidder.Sub.Id)
					subs = db.GetSubs("name=\"" + id + "\"")
					if len(subs) == 0 {
						http.Redirect(w, r, "/static/404.html", http.StatusFound)
						return
					}
					subtidder.Sub = subs[0]
				}
			}
			if true {
				var banner string
				file, header, err := r.FormFile("banner_input")
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
					banner = "data:" + header.Header.Get("Content-Type") + ";base64," // file to base64

					bytes, err := io.ReadAll(file)
					if err != nil {
						Util.Error(err)
					}

					banner += base64.StdEncoding.EncodeToString(bytes)
				}

				if banner != "" {
					db.UpdateSubtidder("banner", banner, subtidder.Sub.Id)
					subs = db.GetSubs("name=\"" + id + "\"")
					if len(subs) == 0 {
						http.Redirect(w, r, "/static/404.html", http.StatusFound)
						return
					}
					subtidder.Sub = subs[0]
				}
			}
		}

		subtidder.Posts = db.GenerateSubTidderFeed(IAM, subtidder.Sub.Id)
		subtidder.SubscribedUser = db.GenerateAccountSubscribed(subtidder.Sub.Id)

		type FetchQuery struct {
			IdPost    int `json:"id_post"`
			Score     int `json:"score"`
			IdAccount int `json:"id_account_subscribing"`
			IdSubject int `json:"id_subject_to_subscribe"`

			Info string `json:"info"`

			BannedChanges string `json:"banned_user_changes"`

			RoleToCreate string `json:"role_to_create"`
			RoleToDelete string `json:"role_to_delete"`
			RoleToUpdate string `json:"role_to_update"`

			RoleAtribution string `json:"role_atribution_changes"`
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

		if fetchQuery.Info != "" {
			db.EditInfo(fetchQuery.Info, subtidder.Sub.Id)
		}

		if fetchQuery.BannedChanges != "" {
			db.ChangeBanned(subtidder.Sub.Id, fetchQuery.BannedChanges)
		}

		if fetchQuery.RoleToCreate != "" {
			db.CreateRole(subtidder.Sub.Id, fetchQuery.RoleToCreate)
		}

		if fetchQuery.RoleToDelete != "" {
			db.DeleteRole(subtidder.Sub.Id, fetchQuery.RoleToDelete)
		}

		if fetchQuery.RoleToUpdate != "" {
			db.UpdateRole(subtidder.Sub.Id, fetchQuery.RoleToUpdate)
		}

		if fetchQuery.RoleAtribution != "" {
			db.ChangeRoleAtribution(subtidder.Sub.Id, fetchQuery.RoleAtribution)
		}

		subtidder.Subscribed = db.IsSubscribeTo(IAM, subtidder.Sub.Id)
		viewData.SubtidderVD = subtidder

		fmt.Println("RoleToCreate : ", fetchQuery.RoleToCreate)
		fmt.Println("RoleToDelete : ", fetchQuery.RoleToDelete)
		fmt.Println("RoleToUpdate : ", fetchQuery.RoleToUpdate)

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
			viewData.Errors.Signup = db.CreateAccount(pseudo, email, password, birthdate, verifpassword)
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
	viewData.Page = "signin"

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
				viewData.Account = connectedUsr
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
	viewData.Page = "post"

	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		IAM := testConnection(r, &viewData, db)
		postVD := SQL.PostVD{}

		id := getSuffix(r)
		_, err := strconv.Atoi(id)
		if id == "!err404" || err != nil {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}

		posts := db.GetPosts("id_post=" + id)
		if len(posts) == 0 {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}
		post := posts[0]

		if r.URL.Path == "/post/default.png" {
			Util.Warning("Laurie fait chier : default.png")
			return
		}

		postVD.Post = db.MakeDisplayablePost(post, IAM)

		subs := db.GetSubs("id_subject=" + strconv.Itoa(post.IdSubject))
		if len(subs) == 0 {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}
		postVD.Subtidder = subs[0]

		comments := db.GetComments("id_post=" + id)
		postVD.Comments, _ = db.MakeDisplayableComments(comments, IAM)

		viewData.PostVD = postVD

		// BASIC NEW COMMENTS //
		if r.Method == "POST" {
			if r.FormValue("comment_content") != "" {
				db.CreateComment(r.FormValue("comment_content"), IAM, id, "-1")
				http.Redirect(w, r, "/post/"+id, http.StatusFound)
			}
		}

		// VOTES //
		type FetchQuery struct {
			IdPost    int `json:"id_post"`
			Score     int `json:"score"`
			IdAccount int `json:"id_account_subscribing"`
			IdSubject int `json:"id_subject_to_subscribe"`

			IdResponseTo int    `json:"id_response_to"`
			Content      string `json:"content"`

			IdComment    int `json:"id_comment"`
			ScoreComment int `json:"score_comment"`
		}
		fetchQuery := &FetchQuery{}
		json.NewDecoder(r.Body).Decode(fetchQuery)

		if fetchQuery.IdPost != 0 && fetchQuery.Score != 0 && IAM != -1 { // vote for post
			db.Vote(fetchQuery.IdPost, fetchQuery.Score, IAM)
		}
		if fetchQuery.IdComment != 0 && fetchQuery.ScoreComment != 0 && IAM != -1 { // vote for comment
			db.VoteForComment(fetchQuery.IdComment, fetchQuery.ScoreComment, IAM)
		}
		if fetchQuery.IdResponseTo != 0 && fetchQuery.Content != "" && IAM != -1 { // comment a comment
			db.CreateComment(fetchQuery.Content, IAM, id, strconv.Itoa(fetchQuery.IdResponseTo))
		}

		err = callTemplate("post_page", &viewData, w)
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
		var profilePageVD SQL.ProfilePageVD

		name := getSuffix(r)
		if name == "!err404" {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}
		profilePageVD.Account = db.GetAccountByName(name)
		if profilePageVD.Account.Id == -1 {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}

		posts := db.GetPosts("id_author =" + strconv.Itoa(profilePageVD.Account.Id))
		for _, post := range posts {
			profilePageVD.Posts = append(profilePageVD.Posts, db.MakeDisplayablePost(post, IAM))
		}
		profilePageVD.Subtidders = db.GetSubtiddersSubscribed(profilePageVD.Account.Id)

		viewData.ProfilePageVD = profilePageVD

		err := callTemplate("profile_page", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func Error404Handler(db *SQL.SqlServer) {

}

func CGU(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	http.HandleFunc("/cgu", func(w http.ResponseWriter, r *http.Request) {
		testConnection(r, &viewData, db)
		name := getSuffix(r)
		if name == "!err404" {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}

		err := callTemplate("cgu", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}

func Authors_redirect(db *SQL.SqlServer) {
	viewData := SQL.MasterVD{}
	http.HandleFunc("/Authors", func(w http.ResponseWriter, r *http.Request) {
		testConnection(r, &viewData, db)
		name := getSuffix(r)
		if name == "!err404" {
			http.Redirect(w, r, "/static/404.html", http.StatusFound)
			return
		}

		err := callTemplate("authors_redirect", &viewData, w)
		if err != nil {
			Util.Error(err)
		}
		viewData.ClearErrors()
	})
}
