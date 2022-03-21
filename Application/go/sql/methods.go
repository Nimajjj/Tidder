package mySQL

import(
  "fmt"
  "strconv"
)

/*
  /go/sql/methods.go

  Scripts only containing utility methods of SqlServer struct.
*/


/*
  (sqlServ SqlServer) GetAccountById(id int) Accounts

  Get Accounts from an id into the database.
*/
func (sqlServ SqlServer) GetAccountById(id int) Accounts {
  fmt.Println("Requesting account", strconv.Itoa(id), "from Accounts table.")
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
  fmt.Println("Request completed.")
  return account
}
