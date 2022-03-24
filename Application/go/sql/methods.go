package mySQL

import (
	"fmt"
	Util "github.com/Nimajjj/Tidder/go/utility"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

/*
  /go/sql/methods.go

  Script only containing utility methods of SqlServer struct.
*/

/*
  (sqlServ SqlServer) GetAccountById(id int) Accounts

  Get Accounts from an id into the database.
*/
func (sqlServ SqlServer) GetAccountById(id int) Accounts {
	var account Accounts
	query := "SELECT * FROM accounts WHERE id_account = " + strconv.Itoa(id)
	Util.Log("Executing following query :")
	Util.Log(query)
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
	if err != nil {
		Util.Error(err)
	}
	return account
}

/*
  (sqlServ SqlServer) CreateSub(subName string, ownerId int, nsfwInput bool)

  Provide a name, an id and a bool to create a new subtidder

  to do :
    check nsfw
*/
func (sqlServ SqlServer) CreateSub(subName string, ownerId int, nsfwInput bool) {
	if len(subName) >= 25 || ownerId < 1 {
		return
	}

	nsfw := 0
	if nsfwInput {
		nsfw = 1
	}

	// test if sub name is already taken
	query := "name=\"" + subName + "\""
	if len(sqlServ.GetSubs(query)) != 0 {
		Util.Log("Creating sub failed : sub name <" + subName + "> already taken.")
		return
	}

	query = "INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("
	query += "\"" + subName + "\", \"default.png\", "
	query += strconv.Itoa(ownerId) + ", " + strconv.Itoa(nsfw) + ")"
	sqlServ.executeQuery(query)
}

/*
  (sqlServ SqlServer) GetSubs(conditions string) []Subject

  Function returning a list a all sub (responding to certain conditions)
*/
func (sqlServ SqlServer) GetSubs(conditions string) []Subject {
	Util.Log("Executing following query :")
	query := "SELECT * FROM subjects "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Log(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}

	result := []Subject{}
	for rows.Next() {
		var id int
		var name string
		var pp string
		var nsfw bool
		var id_owner int
		if err2 := rows.Scan(
			&id,
			&name,
			&pp,
			&id_owner,
			&nsfw,
		); err2 != nil {
			Util.Error(err2)
		}
		sub := Subject{id, name, pp, nsfw, id_owner}
		result = append(result, sub)
	}

	return result
}


func HashPassword(password string) string {
  var passwordBytes = []byte(password)
  hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	if err != nil { Util.Log.Error(err) }
	Util.Log("Password Hash : " + hashedPasswordBytes)

	return hashedPasswordBytes
}


func (sqlServ SqlServer) CreateAccount(name string, email string, Password string, Birthdate string) {
	currentTime := time.Now()
	query := "INSERT INTO accounts (name, email, hashed_password, birth_date , creation_date , karma , profile_picture) VALUES ("
	query += "\"" + name + "\","
	query += "\"" + email + "\","
	query += "\"" + HashPassword(Password) + "\","
	query += "\"" + Birthdate + "\","
	query += "\"" + currentTime.Format("2006-01-02") + "\","
	query += " 0, \"Default.png\" )"
	sqlServ.executeQuery(query)
}
