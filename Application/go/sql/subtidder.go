package mySQL

import (
	Util "github.com/Nimajjj/Tidder/go/utility"
	"strconv"
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
				return "Creating sub failed : sub name <" + subName + "> contains forbidden characters."
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
		return "Creating sub failed : sub name <" + subName + "> already taken."
	}

	query = "INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("
	query += "\"" + subName + "\", \"default_pp.png\", "
	query += strconv.Itoa(ownerId) + ", " + strconv.Itoa(nsfw) + ")"
	Util.Query(query)
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
	Util.Query(query)
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
		result = append(result, sub)
	}

	return result
}


func (sqlServ SqlServer) GetNumberOfSubscriber(subject_id int) int {
	query := "SELECT * FROM subscribe_to_subject WHERE id_subject=" + strconv.Itoa(subject_id)
	Util.Query(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {	Util.Error(err)	}
	result := 0
	for rows.Next() {
		result += 1
	} 
	return result
}

func (sqlServ SqlServer) GetSubtiddersSubscribed(account_id int) []Subject {
	subscribedID :=  []int{} 
	query := "SELECT id_subject FROM tidder.subscribe_to_subject WHERE id_account =" + strconv.Itoa(account_id)
	Util.Query(query)

	rows, err := sqlServ.db.Query(query)
	if err != nil {	Util.Error(err)	}
	
	for rows.Next() {
		var id int
		if err2 := rows.Scan(
			&id,
		); err2 != nil {
			Util.Error(err2)
		}
		subscribedID = append(subscribedID, id)
	} 

	if (len(subscribedID) == 0) {
		return []Subject{}
	}

	result := []Subject{}
	query = "SELECT * FROM tidder.subjects WHERE"
	for i, id := range subscribedID {
		if (i != 0) {
			query += " OR"
		}
		query += " id_subject = " + strconv.Itoa(id)
	}


	Util.Query(query)
	rows2, err2 := sqlServ.db.Query(query)
	if err2 != nil { Util.Error(err2) }
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