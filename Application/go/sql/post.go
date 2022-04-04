package mySQL

import (
	Util "github.com/Nimajjj/Tidder/go/utility"
	"strconv"
)

/*
  /go/sql/posts.go

  Script only containing post related methods of SqlServer struct.
*/

/*
  (sqlServ SqlServer) GetPosts(conditions string) []Posts

  Function returning a list a all posts (responding to certain conditions)
*/
func (sqlServ SqlServer) GetPosts(conditions string) []Posts {
	query := "SELECT * FROM posts "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Query(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}

	result := []Posts{}
	for rows.Next() {
		var id int
		var title string
		var media_url string
		var content string
		var creation_date string
		var upvotes int
		var downvotes int
		var nsfw bool
		var redacted bool
		var pinned bool
		var id_subject int
		var id_author int
		if err2 := rows.Scan(
			&id,
			&title,
			&media_url,
			&content,
			&creation_date,
			&upvotes,
			&downvotes,
			&nsfw,
			&redacted,
			&pinned,
			&id_subject,
			&id_author,
		); err2 != nil {
			Util.Error(err2)
		}
		numberOfComments := len(sqlServ.GetComments("id_post = " + strconv.Itoa(id) + " ORDER BY creation_date DESC"))
		post := Posts{id, title, media_url, content, creation_date, upvotes, downvotes, nsfw, redacted, pinned, id_subject, id_author, upvotes - downvotes, numberOfComments}
		result = append(result, post)
	}

	return result
}

func (sqlServ SqlServer) CreatePost(postName string, ownerId int) {
	if len(postName) >= 125 || ownerId < 1 {
		Util.Warning("Creating post failed : post name <" + postName + "> is longer than 125 characters.")
		return
	}
	forbiddenChar := []string{"`", "|", "*", "&", "@", "~", "^", "{", "}"}
	for _, char := range postName {
		for _, forbidden := range forbiddenChar {
			if string(char) == forbidden {
				Util.Warning("Creating sub failed : sub name <" + postName + "> contains forbidden characters.")
				return
			}
		}
	}

	// test if sub name is already taken
	query := "name=\"" + postName + "\""
	query = "INSERT INTO `subjects` (name, profile_picture, id_owner, nsfw) VALUES ("
	query += "\"" + postName + "\", \"default_pp.png\", "
	query += strconv.Itoa(ownerId) + ", " + ")"
	Util.Query(query)
	sqlServ.executeQuery(query)
}
