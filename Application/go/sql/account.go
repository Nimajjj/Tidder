package mySQL

import (
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	Util "github.com/Nimajjj/Tidder/go/utility"
)

func (sqlServ SqlServer) GetAccount(conditions string) []Accounts {
	query := "SELECT * FROM accounts "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Query("GetAccount", query)
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
	Util.Query("GetAccountById", query)
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

func (sqlServ SqlServer) GetAccountByName(name string) Accounts {
	var account Accounts
	query := "SELECT * FROM accounts WHERE name = " + name
	Util.Query("GetAccountByName", query)
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

func (sqlServ SqlServer) CreateAccount(name string, email string, Password string, Birthdate string, Verif_password string) string {
	error := ""
	if Verif_password != Password {
		error += "Password and verification password are different"
		Util.Warning("Password and verification password are different : " + Password + " != " + Verif_password)
		return error
	}
	if name == "" || email == "" || Password == "" || Birthdate == "" {
		error += "Please complete all fields"
		Util.Warning("User try to create an account with empty fields")
		return error
	}

	currentTime := time.Now()
	query := "name=\"" + name + "\""
	if len(sqlServ.GetAccount(query)) != 0 {
		error += "This pseudoname is already used"
		Util.Warning("Creating account failed : name :  <" + name + "> already taken.")
		return error
	}
	query = "email=\"" + email + "\""
	if len(sqlServ.GetAccount(query)) != 0 {
		error += "This email is already used"
		Util.Warning("Creating account failed : email : <" + email + "> already taken.")
		return error
	}

	query = "INSERT INTO accounts (name, email, hashed_password, birth_date , creation_date , karma , profile_picture) VALUES ("
	query += "\"" + name + "\","
	query += "\"" + email + "\","
	query += "\"" + HashPassword(Password) + "\","
	query += "\"" + Birthdate + "\","
	query += "\"" + currentTime.Format("2006-01-02") + "\","
	query += " 0, \"default.png\")"
	sqlServ.executeQuery(query)
	return error
}

func (sqlServ SqlServer) IsSubscribeTo(idAccount int, idSubject int) bool {
	alreadySubscribed := false
	query := "SELECT * FROM subscribe_to_subject WHERE id_account = " + strconv.Itoa(idAccount) + " AND id_subject = " + strconv.Itoa(idSubject)
	Util.Query("IsSubscribeTo", query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	for rows.Next() {
		alreadySubscribed = true
		break
	}
	return alreadySubscribed
}

func (sqlServ SqlServer) SubscribeToSubject(idAccount int, idSubject int) {
	if !sqlServ.IsSubscribeTo(idAccount, idSubject) {
		query := "INSERT INTO subscribe_to_subject (id_account, id_subject) VALUES ("
		query += strconv.Itoa(idAccount) + ","
		query += strconv.Itoa(idSubject) + ")"
		sqlServ.executeQuery(query)
		Util.Log("User id " + strconv.Itoa(idAccount) + " subscribed to subject id " + strconv.Itoa(idSubject))
	} else {
		query := "DELETE FROM subscribe_to_subject WHERE id_account ="
		query += strconv.Itoa(idAccount) + " AND id_subject = " + strconv.Itoa(idSubject)
		sqlServ.executeQuery(query)
		Util.Log("User id " + strconv.Itoa(idAccount) + " unsubscribed from subject id " + strconv.Itoa(idSubject))
	}
}
