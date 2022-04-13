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

func (sqlServ SqlServer) CreateAccount(name string, email string, Password string, Birthdate string, studentId string, Verif_password string) string {
	error := ""
	if Verif_password != Password {
		error += "Les Mots de passes ne sont pas identiques"
		Util.Warning("Les mots de passes : " + Password + "et" + Verif_password + "ne sont pas identiques")
		return error
	}
	currentTime := time.Now()
	query := "name=\"" + name + "\""
	if len(sqlServ.GetAccount(query)) != 0 {
		error += "Rentrez un pseudo non utilisé"
		Util.Warning("Creating account failed : name :  <" + name + "> already taken.")
		return error
	}
	query = "email=\"" + email + "\""
	if len(sqlServ.GetAccount(query)) != 0 {
		error += "Rentrez un email non utilisé"
		Util.Warning("Creating account failed : email : <" + email + "> already taken.")
		return error
	}
	query = "student_id=\"" + studentId + "\""
	if len(sqlServ.GetAccount(query)) != 0 {
		error += "Rentrez un identifiant ynov non utilisé"
		Util.Warning("Creating account failed : studenId : <" + studentId + "> already taken.")
		return error
	}
	query = "INSERT INTO accounts (name, email, hashed_password, birth_date , creation_date , karma , profile_picture, student_id) VALUES ("
	query += "\"" + name + "\","
	query += "\"" + email + "\","
	query += "\"" + HashPassword(Password) + "\","
	query += "\"" + Birthdate + "\","
	query += "\"" + currentTime.Format("2006-01-02") + "\","
	query += " 0, \"Default.png\","
	query += "\"" + studentId + "\")"
	sqlServ.executeQuery(query)
	return error
}

func (sqlServ SqlServer) SubscribeToSubject(idAccount int, idSubject int) {
	alreadySubscribed := false

	query := "SELECT * FROM subscribe_to_subject WHERE id_account = " + strconv.Itoa(idAccount) + " AND id_subject = " + strconv.Itoa(idSubject)
	Util.Query(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	for rows.Next() {
		alreadySubscribed = true
		break
	}

	if !alreadySubscribed {
		query = "INSERT INTO subscribe_to_subject (id_account, id_subject) VALUES ("
		query += strconv.Itoa(idAccount) + ","
		query += strconv.Itoa(idSubject) + ")"
		Util.Query(query)
		sqlServ.executeQuery(query)
		Util.Log("User id " + strconv.Itoa(idAccount) + " subscribed to subject id " + strconv.Itoa(idSubject))
	} else {
		query = "DELETE FROM subscribe_to_subject WHERE id_account ="
		query += strconv.Itoa(idAccount) + " AND id_subject = " + strconv.Itoa(idSubject)
		Util.Query(query)
		sqlServ.executeQuery(query)
		Util.Log("User id " + strconv.Itoa(idAccount) + " unsubscribed from subject id " + strconv.Itoa(idSubject))
	}
}

func (sqlServ SqlServer) Connection(idAccount int, name string, password string, Isconnected bool) {
	query := "SELECT * FROM accounts WHERE id_account = " + strconv.Itoa(idAccount)
	firstquery := "SELECT PASSWORD FROM accounts WHERE id_account = " + strconv.Itoa(idAccount)
	if password == firstquery {
		sqlServ.executeQuery(query)
		Isconnected = true
	}

}
