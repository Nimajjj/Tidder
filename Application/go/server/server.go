package server

import (
	SQL "github.com/Nimajjj/Tidder/go/sql"
	Util "github.com/Nimajjj/Tidder/go/utility"
	"html/template"
	"net/http"
)

/*
  /go/server/server.go

  All files path must be relative to the main.go script position.
  Only functions starting with a CAPITAL LETTER can be access from the main.go script
*/

/*
  Run()

  ONLY function which can be accessed from the main.go script
  Entry door for the go web server.
*/

func Run(DatabaseIp string) {
	Util.Log("Tidder Inc © 2022. Tous droits réservés")
	Util.Log("Starting server : http://localhost:80")

	initStaticFolders()
	launchServer(DatabaseIp)
}

/*
  initStaticFolders()

  All statics folders which will be used in html/css/js files must be declared here.
*/
func initStaticFolders() {

	cssFolder := http.FileServer(http.Dir("./style"))
	imgFolder := http.FileServer(http.Dir("./images"))
	jsFolder := http.FileServer(http.Dir("./scripts"))
	http.Handle("/style/", http.StripPrefix("/style/", cssFolder))
	http.Handle("/images/", http.StripPrefix("/images/", imgFolder))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", jsFolder))
}

/*
  launchServer()

  Main function which run the server
  For the moment the template & the http.ListenAndServe(":80", nil) are here.

  To do :
    -create an individual function for each template.
*/

func launchServer(DatabaseIp string) {
	var db SQL.SqlServer
	db.Connect(DatabaseIp)
	defer db.Close()

	IndexHandler(&db)
	SubtidderHandler(&db)
	SearchHandler(&db)

	testTpl := template.Must(template.ParseFiles("./test/index2.html"))
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		type Error struct {
			Error string
		}
		err := Error{""}
		if r.Method == http.MethodPost {
			pseudo := r.FormValue("pseudo_input")
			email := r.FormValue("email_input")
			password := r.FormValue("password_input")
			birthdate := r.FormValue("birthdate_input")
			studentId := r.FormValue("id_input")
			if pseudo != "" && email != "" && password != "" && birthdate != "" && studentId != "" {
				err.Error = db.CreateAccount(pseudo, email, password, birthdate, studentId)
			} else {
				err.Error = "Rentrez des informations valides"
			}
		}
		testTpl.Execute(w, err)
	})

	http.ListenAndServe(":80", nil)
}
