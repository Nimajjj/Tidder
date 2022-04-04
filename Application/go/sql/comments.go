package mySQL

import (
	Util "github.com/Nimajjj/Tidder/go/utility"
)

func (sqlServ SqlServer) GetComments(conditions string) []Comments {
	query := "SELECT * FROM comments "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Query(query)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}

	result := []Comments{}
	for rows.Next() {
		var id int
		var content string
		var creation_date string
		var upvotes int
		var downvotes int
		var redacted bool
		var id_author int
		var response_to_id int
		var id_post int
		if err2 := rows.Scan(
			&id,
			&content,
			&creation_date,
			&upvotes,
			&downvotes,
			&redacted,
			&id_author,
			&response_to_id,
			&id_post,
		); err2 != nil {
			Util.Error(err2)
		}
		comment := Comments{id, content, creation_date, upvotes, downvotes, redacted, id_author, response_to_id, id_post, upvotes - downvotes}
		result = append(result, comment)
	}

	return result
}

func (sqlServ SqlServer) Creatingcomment() {

}