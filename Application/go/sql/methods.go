package mySQL

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
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
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Request completed.")
	return account
}

// /!\ Not finished yet /!\ //
// INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("Dank Meme", "default.png", 1, 0)
func (sqlServ SqlServer) CreateSub(subName string, ownerId int, nsfwInput bool) {
	if len(subName) >= 25 || ownerId < 1 {
		return
	}
	nsfw := 0
	if nsfwInput {
		nsfw = 1
	}
	query := "INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("
	query += "\"" + subName + "\", \"default.png\", "
	query += strconv.Itoa(ownerId) + ", " + strconv.Itoa(nsfw)
	sqlServ.executeQuery(query)
}

func (sqlServ SqlServer) executeQuery(query string) {
	fmt.Println("Executing following query :")
	fmt.Println(query)
	_, err := sqlServ.db.Query(query)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("Query successfully executed.")
}

func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func HashPassword(password string) string {

	var passwordres, err = hashPassword(password)
	if err != nil {
		println(fmt.Println("Error hashing password"))
		return passwordres
	}
	fmt.Println("Password Hash:", passwordres)
	return passwordres
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
