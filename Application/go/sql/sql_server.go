package mySQL

import (
  "strconv"
  "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


type SqlServer struct {
  db *sql.DB
}


func (sqlServ *SqlServer) Connect() {
  db, err := sql.Open("mysql", "root:Tidder123reddit@tcp(127.0.0.1:3306)/tidder")
  if err != nil { panic(err.Error()) }
  sqlServ.db = db
}


func (sqlServ SqlServer) Close() {
  sqlServ.db.Close()
}


func (sqlServ SqlServer) GetDB() *sql.DB {
  return sqlServ.db
}


func (sqlServ SqlServer) GetAccountById(id int) Accounts {
  var account Accounts
  query := "SELECT * FROM accounts WHERE id_account = " + strconv.Itoa(id)
  err := sqlServ.db.QueryRow(query).Scan(
    &account.Id,
    &account.Name,
    &account.Email,
    &account.Password,
    &account.BirthDate,
    &account.CreationDate,
    &account.Karma,
    &account.ProfilePicture,
  )
  if err != nil { panic(err.Error()) }
  return account
}
