package mySQL

import (
	"strconv"

	Util "github.com/Nimajjj/Tidder/go/utility"
)

/*
  (sqlServ SqlServer) CreateSub(subName string, ownerId int, nsfwInput bool)

  Provide a name, an id and a bool to create a new subtidder

  to do :
    check nsfw
*/
func (sqlServ SqlServer) CreateSub(subName string, ownerId int, nsfwInput bool) string {
	if len(subName) >= 25 || ownerId < 1 {
		Util.Warning("Creating sub failed : sub name <" + subName + "> is longer than 24 characters.")
		return "Creating sub failed : sub name <" + subName + "> is longer than 24 characters."
	}
	forbiddenChar := []string{" ", "'", "\"", "`", ";", ",", ".", ":", "!", "?", "\\", "/", "|", "=", "*", "&", "%", "$", "#", "@", "~", "^", "(", ")", "[", "]", "{", "}", "<", ">"}
	for _, char := range subName {
		for _, forbidden := range forbiddenChar {
			if string(char) == forbidden {
				Util.Warning("Creating sub failed : sub name <" + subName + "> contains forbidden characters.")
				return "Subtidder name \"" + subName + "\" contains forbidden characters."
			}
		}
	}

	nsfw := 0
	if nsfwInput {
		nsfw = 1
	}

	// test if sub name is already taken
	query := "name=\"" + subName + "\""
	if len(sqlServ.GetSubs(query)) != 0 {
		Util.Warning("Creating sub failed : sub name <" + subName + "> already taken.")
		return "Subtidder name \"" + subName + "\" already taken."
	}

	query = "INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw, banner) VALUES ("
	query += "\"" + subName + "\", \"\", "
	query += strconv.Itoa(ownerId) + ", " + strconv.Itoa(nsfw) + ", \"\")"
	sqlServ.executeQuery(query)
	return ""
}

/*
  (sqlServ SqlServer) GetSubs(conditions string) []Subject

  Function returning a list a all subjects (responding to certain conditions)
*/
func (sqlServ SqlServer) GetSubs(conditions string) []Subject {
	query := "SELECT * FROM subjects "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Query("GetSubs", query)
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
		var infos string
		var banner string
		if err2 := rows.Scan(
			&id,
			&name,
			&pp,
			&id_owner,
			&nsfw,
			&banner,
			&infos,
		); err2 != nil {
			Util.Error(err2)
		}
		sub := Subject{id, name, pp, nsfw, id_owner, infos, banner}
		if banner == "" {
			sub.Banner = DefaultSubtidderBanner()
		}
		if pp == "" {
			sub.ProfilePicture = DefaultSubtidderPP()
		}
		result = append(result, sub)
	}

	return result
}

func (sqlServ SqlServer) GetNumberOfSubscriber(subject_id int) int {
	query := "SELECT * FROM subscribe_to_subject WHERE id_subject=" + strconv.Itoa(subject_id)
	Util.Query("GetNumberOfSubscriber", query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	result := 0
	for rows.Next() {
		result += 1
	}
	return result
}

// GetSubtiddersSubscribed return array of all subtidder where given account is subscribe
func (sqlServ SqlServer) GetSubtiddersSubscribed(account_id int) []Subject {
	if account_id == -1 {
		return []Subject{}
	}

	subscribedID := []int{}
	query := "SELECT id_subject FROM tidder.subscribe_to_subject WHERE id_account =" + strconv.Itoa(account_id)
	Util.Query("GetSubtiddersSubscribed", query)

	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}

	for rows.Next() {
		var id int
		if err2 := rows.Scan(
			&id,
		); err2 != nil {
			Util.Error(err2)
		}
		subscribedID = append(subscribedID, id)
	}

	if len(subscribedID) == 0 {
		return []Subject{}
	}

	result := []Subject{}
	query = "SELECT * FROM tidder.subjects WHERE"
	for i, id := range subscribedID {
		if i != 0 {
			query += " OR"
		}
		query += " id_subject = " + strconv.Itoa(id)
	}

	Util.Query("GetSubtiddersSubscribed", query)
	rows2, err2 := sqlServ.db.Query(query)
	if err2 != nil {
		Util.Error(err2)
	}
	for rows2.Next() {
		var id int
		var name string
		var pp string
		var nsfw bool
		var id_owner int
		var infos string
		var banner string
		if err2 := rows2.Scan(
			&id,
			&name,
			&pp,
			&id_owner,
			&nsfw,
			&banner,
			&infos,
		); err2 != nil {
			Util.Error(err2)
		}
		sub := Subject{id, name, pp, nsfw, id_owner, infos, banner}
		result = append(result, sub)
	}

	return result
}

func (sqlServ SqlServer) EditInfo(newInfo string, subjectId int) {
	query := "UPDATE subjects SET infos=" + "\"" + newInfo + "\"" + " WHERE id_subject=" + strconv.Itoa(subjectId)
	Util.Query("EditInfo", query)
	sqlServ.executeQuery(query)
}

func (sqlServ SqlServer) UpdateSubtidder(col string, val string, id int) {
	query := "UPDATE subjects SET " + col + "='" + val + "' WHERE id_subject=" + strconv.Itoa(id)
	Util.Query("UpdateSubtidder", query)
	sqlServ.executeQuery(query)
}

func (sqlServ SqlServer) GenerateAccountSubscribed(idSubject int) []AccountSubscribed {
	accountSubscribed := []AccountSubscribed{}

	query := "SELECT * FROM subscribe_to_subject WHERE id_subject=" + strconv.Itoa(idSubject)
	Util.Query("GenerateBannedAccountList", query)
	rows, err := sqlServ.db.Query(query)

	// Get all subscribed account
	IDs := []int{}
	if err != nil {
		Util.Error(err)
	} else {
		for rows.Next() {
			var idAccount int
			var idSubject int
			if err2 := rows.Scan(
				&idAccount,
				&idSubject,
			); err2 != nil {
				Util.Error(err2)
			}
			IDs = append(IDs, idAccount)
		}
	}

	if len(IDs) == 0 {
		return accountSubscribed
	}

	// Get all banned account
	query = "SELECT * FROM accounts WHERE "
	for i, id := range IDs {
		query += "id_account=" + strconv.Itoa(id) + " "
		if i != len(IDs)-1 {
			query += "OR "
		}
	}

	Util.Query("GenerateBannedAccountList", query)
	rows, err = sqlServ.db.Query(query)
	accounts := []Accounts{}

	if err != nil {
		Util.Error(err)
	} else {
		for rows.Next() {
			var id int
			var name string
			var email string
			var password string
			var birthdate string
			var creationDate string
			var karma int
			var pp string
			var studentID string
			if err2 := rows.Scan(
				&id,
				&name,
				&email,
				&password,
				&birthdate,
				&creationDate,
				&karma,
				&pp,
				&studentID,
			); err2 != nil {
				Util.Error(err2)
			}
			accounts = append(accounts, Accounts{id, name, email, password, birthdate, creationDate, karma, pp, studentID})
		}
	}

	query = "SELECT * FROM is_ban WHERE id_subject=" + strconv.Itoa(idSubject)
	Util.Query("GenerateBannedAccountList", query)
	rows, err = sqlServ.db.Query(query)
	idBanned := []int{}

	if err != nil {
		Util.Error(err)
	} else {
		for rows.Next() {
			var id int
			var osef string
			if err2 := rows.Scan(
				&id,
				&osef,
			); err2 != nil {
				Util.Error(err2)
			}
			idBanned = append(idBanned, id)
		}
	}

	for _, account := range accounts {
		var formatAccount AccountSubscribed
		formatAccount.Account = account
		formatAccount.Banned = false

		defaultRole := SubjectRoles{}

		if sqlServ.RowExists("has_subject_role", "id_account="+strconv.Itoa(account.Id)+" AND id_subject="+strconv.Itoa(idSubject)) {
			query = "SELECT * FROM has_subject_role WHERE id_account=" + strconv.Itoa(account.Id) + " AND id_subject=" + strconv.Itoa(idSubject)
			Util.Query("GenerateBannedAccountList", query)
			rows, err = sqlServ.db.Query(query)
			if err != nil {
				Util.Error(err)
			}
			for rows.Next() {
				var idAccount int
				var idSubject int
				var idRole int
				if err2 := rows.Scan(
					&idAccount,
					&idSubject,
					&idRole,
				); err2 != nil {
					Util.Error(err2)
				}
				defaultRole.Id = idRole
				break
			}
		} else {
			defaultRole.Id = -1
		}

		formatAccount.Role = defaultRole

		for _, id := range idBanned {
			if id == account.Id {
				formatAccount.Banned = true
				break
			}
		}

		accountSubscribed = append(accountSubscribed, formatAccount)
	}

	return accountSubscribed
}

func (sqlServ SqlServer) ChangeBanned(idSubject int, unformatedChanges string) {
	changes := []int{}
	n := ""
	for _, char := range unformatedChanges {
		if char != ';' {
			n += string(char)
		} else {
			id, _ := strconv.Atoi(n)
			changes = append(changes, id)
			n = ""
		}
	}

	query := "SELECT * FROM is_ban WHERE id_subject=" + strconv.Itoa(idSubject) + " AND ("

	for i, id := range changes {
		query += "id_account=" + strconv.Itoa(id)
		if i != len(changes)-1 {
			query += " OR "
		}
	}
	query += " )"
	Util.Query("ChangeBanned", query)
	rows, err := sqlServ.db.Query(query)
	toUnban := []int{}

	if err != nil {
		Util.Error(err)
	} else {
		for rows.Next() {
			var id int
			var osef string
			if err2 := rows.Scan(
				&id,
				&osef,
			); err2 != nil {
				Util.Error(err2)
			}
			toUnban = append(toUnban, id)
		}
	}

	toBan := []int{}
	for _, id := range changes {
		ban := true
		for _, idd := range toUnban {
			if id == idd {
				ban = false
				break
			}
		}
		if ban {
			toBan = append(toBan, id)
		}
	}

	for _, id := range toUnban {
		query = "DELETE FROM is_ban WHERE id_subject=" + strconv.Itoa(idSubject) + " AND id_account=" + strconv.Itoa(id)
		sqlServ.executeQuery(query)
	}

	for _, id := range toBan {
		query = "INSERT INTO is_ban (id_subject, id_account) VALUES (" + strconv.Itoa(idSubject) + ", " + strconv.Itoa(id) + ")"
		sqlServ.executeQuery(query)
	}

}
