package mySQL

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  Util "github.com/Nimajjj/Tidder/go/utility"
)

/*
  /go/sql/sql_server.go

  SqlServer is a struct allowing to use the distant database.
  All main methods from the SqlServer struct goes here.

  Exemple :
    var db SQL.SqlServer
    db.Connect()
    defer db.Close()
    account := db.GetAccountById(1)
  End

  username :  root
  password :  Tidder123reddit
  ip :        127.0.0.1
  port :      3306
  db :        tidder
*/

type SqlServer struct {
  db *sql.DB
}

/*
  (sqlServ *SqlServer) Connect()

  Connect structure to the distant mySql database.
  First method to use SqlServer
*/
func (sqlServ *SqlServer) Connect(ip string) {
  Util.Log("Connecting to @tcp(" + ip + ")/tidder ...")
  db, err := sql.Open("mysql", "root:Tidder123reddit@tcp(" + ip + ")/tidder")
  if err != nil { Util.Error(err) }
  sqlServ.db = db
  Util.Log("Connection completed.")
}

/*
  (sqlServ SqlServer) Close()

  Close connection with the distant mySql database.
  Second method to use : MUST USE `defer`
*/
func (sqlServ SqlServer) Close() {
  fmt.Println("Closing @tcp(127.0.0.1:3306)/tidder connection.")
  sqlServ.db.Close()
}

/*
  (sqlServ SqlServer) GetDB() *sql.DB

  Function allowing to get the db var from the sql library
*/
func (sqlServ SqlServer) GetDB() *sql.DB {
  return sqlServ.db
}

/*
  (sqlServ SqlServer) executeQuery(query string)

  Execute brut query
  Private SqlServer method
  To use ONLY if you already made all verifications to avoid SQL error
*/
func (sqlServ SqlServer) executeQuery(query string) {
  Util.Log("Executing following query :")
  Util.Log(query)
  _, err := sqlServ.db.Query(query)
  if err != nil{
    Util.Error(err)
    return
  }
  Util.Log("Query successfully executed.")
}
