package mySQL

import (
	"time"
	"strconv"
	Util "github.com/Nimajjj/Tidder/go/utility"
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
	if err != nil {	Util.Error(err)	}

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
		post := Posts{id, title, media_url, content, creation_date, upvotes, downvotes, nsfw, redacted, pinned, id_subject, id_author, upvotes-downvotes, numberOfComments}
		result = append(result, post)
	}

	return result
}


func (sqlServ SqlServer) CreatePost(title string, media_url string, content string,  nsfw_input bool, id_subject string, id_author int) {
	if len(title) >= 125 || id_author < 1 {
		Util.Warning("Creating post failed : post name <" + title + "> is longer than 125 characters.")
		return
	}
	forbiddenChar := []string{"`", "|", "*", "&", "@", "~", "^", "{", "}"}
	for _, char := range title {
		for _, forbidden := range forbiddenChar {
			if string(char) == forbidden {
				Util.Warning("Creating sub failed : sub name <" + title + "> contains forbidden characters.")
				return
			}
		}
	}
	if id_subject == "-1" {
		Util.Warning("Creating sub failed : please choose a correct subtidder.")
		return
	}

	if content == "" {
		Util.Warning("Creating sub failed : content cannot be empty.")
		return
	}

	if title == "" {
		Util.Warning("Creating sub failed : content cannot be empty.")
		return
	}



	nsfw := "0"
	if nsfw_input == true {
		nsfw = "1"
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := "INSERT INTO posts (title, media_url, content, creation_date, upvotes, downvotes, nsfw, redacted, pinned, id_subject, id_author) VALUES ("
	query += "\"" + title + "\","
	query += "\"" + media_url + "\","
	query += "\"" + content + "\","
	query += "\"" + currentTime + "\","
	query += " 0, 0,"
	query += "\"" + nsfw + "\","
	query += " 0, 0,"
	query += "\"" + id_subject + "\","
	query += "\"" + strconv.Itoa(id_author) + "\")"

	sqlServ.executeQuery(query)
}
