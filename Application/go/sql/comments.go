package mySQL

import (
	"html/template"
	"strconv"
	"time"

	Util "github.com/Nimajjj/Tidder/go/utility"
)

func (sqlServ SqlServer) GetComments(conditions string) []Comments {
	query := "SELECT * FROM comments "
	if conditions != "" {
		query += "WHERE " + conditions
	}
	Util.Query("GetComments", query)
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

func (sqlServ SqlServer) CreateComment(content string, idUser int, idPost string) {
	if content == "" {
		return
	}
	if idUser <= 0 {
		return
	}
	if idPost == "" {
		return
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := "INSERT INTO comments (content, creation_date, id_author, id_post) VALUES ("
	query += "'" + content + "', "
	query += "'" + currentTime + "', "
	query += "'" + strconv.Itoa(idUser) + "', "
	query += "'" + idPost + "')"

	sqlServ.executeQuery(query)
}

func (sqlServ SqlServer) MakeDisplayableComments(comments []Comments) ([]DisplayableComment, []int) {
	res := []DisplayableComment{}

	alreadyMade := []int{}

	for _, c := range comments {
		retry := false
		for _, a := range alreadyMade {
			if a == c.Id {
				retry = true
				break
			}
		}
		if retry {
			continue
		}

		displayableComment := DisplayableComment{}
		displayableComment.Id = c.Id
		displayableComment.Content = c.Content
		displayableComment.CreationDate = c.CreationDate
		displayableComment.Upvotes = c.Upvotes
		displayableComment.Downvotes = c.Downvotes
		displayableComment.Redacted = c.Redacted
		displayableComment.IdPost = c.IdPost
		displayableComment.Score = c.Upvotes - c.Downvotes

		author := sqlServ.GetAccountById(c.IdAuthor)
		displayableComment.AuthorName = author.Name
		displayableComment.AuthorPP = author.ProfilePicture
		if displayableComment.AuthorPP == "" {
			displayableComment.AuthorPP = DefaultPP()
		}

		responses := sqlServ.GetComments("response_to_id = " + strconv.Itoa(c.Id))
		if len(responses) > 0 {
			toAppend := []int{}
			displayableComment.Response, toAppend = sqlServ.MakeDisplayableComments(responses)
			for _, r := range toAppend {
				alreadyMade = append(alreadyMade, r)
			}
		}

		res = append(res, displayableComment)

		alreadyMade = append(alreadyMade, c.Id)
	}

	return res, alreadyMade
}

func (c DisplayableComment) PrintComment() template.HTML {

	html := `<div class="comment">
            <div class="header">
                <img src="` + c.AuthorPP + `" alt="author pp">
                <a href="` + c.AuthorName + `">u/` + c.AuthorName + `</a>
                <p>` + c.CreationDate + `</p>
            </div>
            <div class="footer">
                <div class="bar_container">
                    <div class="bar"></div>
                </div>

                <div class="content">
                    <div class="comment_text">` + c.Content + `</div>
                    <div class="content_footer">
                        <img class="post_icon vote_bt upvote_bt" state="active" src="../images/global/empty_upvote.png" alt=" upvote">
                        <p>` + strconv.Itoa(c.Score) + `</p>
                        <img class="post_icon vote_bt downvote_bt" state="empty"  src="../images/global/empty_downvote.png" alt=" downvote">
                        <p>Answer</p>
                        <p>...</p>
                    </div>
	`
	if len(c.Response) > 0 {
		for _, r := range c.Response {
			html += string(r.PrintComment())
		}
	}
	html += `</div>
            </div>
        </div>
	`

	return template.HTML(html)
}
