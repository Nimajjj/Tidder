package mySQL

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	Util "github.com/Nimajjj/Tidder/go/utility"
	"golang.org/x/crypto/bcrypt"
)

/* to do:
- clear cookies
- clear session
- disconnect option
*/

func (sqlServ SqlServer) GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	token := hex.EncodeToString(b)
	query := "SELECT 'id_session' FROM sessions WHERE id_session = '" + token + "'"
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	toBreak := false
	for rows.Next() {
		toBreak = true
		break
	}
	if toBreak {
		Util.Log("Token already used")
		token = sqlServ.GenerateSecureToken(length)
	}

	return token
}

func (sqlServ SqlServer) TryToConnectUser(usr string, psw string) ([]Accounts, string) {
	query := "name = '" + usr + "'"
	account := sqlServ.GetAccount(query)
	if len(account) == 0 {
		Util.Warning("No account found")
		return []Accounts{}, ""
	}

	err := bcrypt.CompareHashAndPassword([]byte(account[0].Password), []byte(psw))
	if err != nil {
		Util.Warning("Wrong password")
		return []Accounts{}, ""
	}

	token := sqlServ.GenerateSecureToken(10)
	Util.Log("New token : " + token)

	query = "INSERT INTO sessions (id_session, id_account, creation_date) VALUES ('" + token + "', " + strconv.Itoa(account[0].Id) + ", '" + time.Now().Format("2006-01-02 15:04:05") + "')"
	sqlServ.executeQuery(query)
	return account, token
}

func (sqlServ SqlServer) GetAccountFromSession(sessionId string) Accounts {
	query := "SELECT id_account FROM sessions WHERE id_session = '" + sessionId + "'"
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		return sqlServ.GetAccount("id_account = " + strconv.Itoa(id))[0]
	}
	return Accounts{}
}
