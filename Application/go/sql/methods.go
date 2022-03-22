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

// /!\ Not finished yet /!\ //
// INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("Dank Meme", "default.png", 1, 0)
func (sqlServ SqlServer) CreateSub(subName string, ownerId int, nsfwInput bool) {
  if len(subName) >= 25 || ownerId < 1 { return }

  nsfw := 0
  if nsfwInput { nsfw = 1 }

  query := "INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("
  query += "\"" + subName + "\", \"default.png\", "
  query += strconv.Itoa(ownerId) + ", " + strconv.Itoa(nsfw)
  sqlServ.executeQuery(query)
}


func (sqlServ SqlServer) executeQuery(query string) {
  _, err := sqlServ.db.Query(query)
  if err != nil{ panic(err) }
}
