package mySQL

import (
	Util "github.com/Nimajjj/Tidder/go/utility"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func (sqlServ SqlServer) GetAccount(conditions string) []Accounts {
	query := "SELECT * FROM accounts "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Query(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}

	result := []Accounts{}
	for rows.Next() {
		var id int
		var name string
		var email string
		var hashed_password string
		var birth_date string
		var creation_date string
		var karma int
		var profile_picture string
		var student_id string
		if err2 := rows.Scan(
			&id,
			&name,
			&email,
			&hashed_password,
			&birth_date,
			&creation_date,
			&karma,
			&profile_picture,
			&student_id,
		); err2 != nil {
			Util.Error(err2)
		}
		account := Accounts{id, name, email, hashed_password, birth_date, creation_date, karma, profile_picture, student_id}
		result = append(result, account)
	}

	return result
}

/*
  (sqlServ SqlServer) GetAccountById(id int) Accounts

  Get Accounts from an id into the database.
*/
func (sqlServ SqlServer) GetAccountById(id int) Accounts {
	var account Accounts
	query := "SELECT * FROM accounts WHERE id_account = " + strconv.Itoa(id)
	Util.Query(query)
	err := sqlServ.db.QueryRow(query).Scan(
		&account.Id,
		&account.Name,
		&account.Email,
		&account.Password,
		&account.BirthDate,
		&account.CreationDate,
		&account.Karma,
		&account.ProfilePicture,
		&account.StudentId,
	)
	if err != nil {
		Util.Error(err)
	}
	return account
}

func HashPassword(password string) string {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	hashedPassword := string(hashedPasswordBytes)
	if err != nil {
		Util.Error(err)
	}
	Util.Log("Password Hash : " + hashedPassword)
	return hashedPassword
}

func (sqlServ SqlServer) CreateAccount(name string, email string, Password string, Birthdate string, studentId string) {
	currentTime := time.Now()
	query := "name=\"" + name + "\" OR "
	query += "\"" + email + "\" OR "
	query += "\"" + studentId + "\""
	if len(sqlServ.GetAccount(query)) != 0 {
		Util.Warning("Creating sub failed : sub name <" + name + "> already taken.")
		return
	}
	query = "INSERT INTO accounts (name, email, hashed_password, birth_date , creation_date , karma , profile_picture, student_id) VALUES ("
	query += "\"" + name + "\","
	query += "\"" + email + "\","
	query += "\"" + HashPassword(Password) + "\","
	query += "\"" + Birthdate + "\","
	query += "\"" + currentTime.Format("2006-01-02") + "\","
	query += " 0, \"Default.png\","
	query += "\"" + studentId + "\")"
	Util.Query(query)
	sqlServ.executeQuery(query)
}
