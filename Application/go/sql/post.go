package mySQL

import (
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
		post := Posts{id, title, media_url, content, creation_date, upvotes, downvotes, nsfw, redacted, pinned, id_subject, id_author}
		result = append(result, post)
	}

	return result
}
