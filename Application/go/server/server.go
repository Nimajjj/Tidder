package server

import (
  "fmt"
  "html/template"
	"net/http"
)

/*
  All files path must be relative to the main.go script position.
  Only functions starting with a CAPITAL LETTER can be access from the main.go script
*/

func Run()  {
  fmt.Println("\nTidder Inc © 2022. Tous droits réservés")
  fmt.Println("Starting server : http://localhost:80")

  initStaticFolders()
  server()
}


func initStaticFolders() {
  cssFolder := http.FileServer(http.Dir("./style"))
  imgFolder := http.FileServer(http.Dir("./images"))
  http.Handle("/style/", http.StripPrefix("/style/", cssFolder))
  http.Handle("/images/", http.StripPrefix("/images/", imgFolder))
}


func server() {
  indexTpl := template.Must(template.ParseFiles("./pages/index.html"))

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexTpl.Execute(w, nil)
	})

  fmt.Println("Server successfully launched.\n")
  http.ListenAndServe(":80", nil)
}
