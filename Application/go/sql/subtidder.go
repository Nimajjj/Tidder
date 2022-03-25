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

  Function returning a list a all subjects (responding to certain conditions)
*/
func (sqlServ SqlServer) GetSubs(conditions string) []Subject {
	Util.Log("Executing following query :")
	query := "SELECT * FROM subjects "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Log(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {	Util.Error(err)	}

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
