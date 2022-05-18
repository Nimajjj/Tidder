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

func (sqlServ SqlServer) CreateComment(content string, idUser int, idPost string, answerToId string) {
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
	query := "INSERT INTO comments (content, creation_date, id_author, id_post, response_to_id) VALUES ("
	query += "'" + content + "', "
	query += "'" + currentTime + "', "
	query += strconv.Itoa(idUser) + ", "
	query += "'" + idPost + "', "
	query += answerToId + ")"

	sqlServ.executeQuery(query)
}

func (sqlServ SqlServer) MakeDisplayableComments(comments []Comments, IAM int) ([]DisplayableComment, []int) {
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

		vote := 0
		query := "SELECT * FROM vote_comment_to WHERE id_account = " + strconv.Itoa(IAM) + " AND id_comment = " + strconv.Itoa(c.Id)
		Util.Query("MakeDisplayableComments", query)
		rows, err := sqlServ.db.Query(query)
		if err != nil {
			Util.Error(err)
		}
		for rows.Next() {
			var a int
			var b int
			var c int
			if err2 := rows.Scan(
				&a,
				&b,
				&c,
				&vote,
			); err2 != nil {
				Util.Error(err2)
			}
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
		displayableComment.Vote = vote

		author := sqlServ.GetAccountById(c.IdAuthor)
		displayableComment.AuthorName = author.Name
		displayableComment.AuthorPP = author.ProfilePicture
		if displayableComment.AuthorPP == "default" || displayableComment.AuthorPP == "" {
			displayableComment.AuthorPP = DefaultPP()
		}

		responses := sqlServ.GetComments("response_to_id = " + strconv.Itoa(c.Id))
		if len(responses) > 0 {
			toAppend := []int{}
			displayableComment.Response, toAppend = sqlServ.MakeDisplayableComments(responses, IAM)
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
                    <div class="content_footer" id="comment_` + strconv.Itoa(c.Id) + `">`

	if c.Vote == 1 {
		html += `<img class="post_icon vote_bt comments_upvote_bt" state="active" src="../images/global/upvote.png" alt="` + strconv.Itoa(c.Id) + ` upvote">
				<p>` + strconv.Itoa(c.Score) + `</p>
				<img class="post_icon vote_bt comments_downvote_bt" state="empty"  src="../images/global/empty_downvote.png" alt="` + strconv.Itoa(c.Id) + ` downvote">`
	} else if c.Vote == -1 {
		html += `<img class="post_icon vote_bt comments_upvote_bt" state="empty" src="../images/global/empty_upvote.png" alt="` + strconv.Itoa(c.Id) + ` upvote">
				<p>` + strconv.Itoa(c.Score) + `</p>
				<img class="post_icon vote_bt comments_downvote_bt" state="active"  src="../images/global/downvote.png" alt="` + strconv.Itoa(c.Id) + ` downvote">`
	} else {
		html += `<img class="post_icon vote_bt comments_upvote_bt" state="empty" src="../images/global/empty_upvote.png" alt="` + strconv.Itoa(c.Id) + ` upvote">
				<p>` + strconv.Itoa(c.Score) + `</p>
				<img class="post_icon vote_bt comments_downvote_bt" state="empty"  src="../images/global/empty_downvote.png" alt="` + strconv.Itoa(c.Id) + ` downvote">`
	}

	html += `<a onClick="AnswerToComment(` + strconv.Itoa(c.Id) + `)">Answer</a> <p>...</p> </div>`

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

func (sqlServ SqlServer) VoteForComment(idComment int, score int, idAccount int) {
	Util.Log("Vote for comment received: " + strconv.Itoa(idComment) + " " + strconv.Itoa(score))
	query := "UPDATE tidder.comments SET "
	if score < 0 {
		query += "downvotes = downvotes + " + strconv.Itoa(score*(-1))
	} else {
		query += "upvotes = upvotes + " + strconv.Itoa(score)
	}
	query += " WHERE id_comment = " + strconv.Itoa(idComment)
	sqlServ.executeQuery(query)

	type VoteTo struct {
		IdComment int
		IdAccount int
		Vote      int
		Shit      int
	}

	var voteTo VoteTo
	query = "SELECT * FROM vote_comment_to WHERE id_account = " + strconv.Itoa(idAccount) + " AND id_comment = " + strconv.Itoa(idComment)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	res := 0
	for rows.Next() {
		res += 1
		var id_comment int
		var id_account int
		var vote int
		var shit int
		if err2 := rows.Scan(
			&shit,
			&id_comment,
			&id_account,
			&vote,
		); err2 != nil {
			Util.Error(err2)
		}
		voteTo = VoteTo{id_comment, id_account, vote, shit}
	}
	if res == 0 {
		query = "INSERT INTO vote_comment_to (id_comment, id_account, vote) VALUES (" + strconv.Itoa(idComment) + ", " + strconv.Itoa(idAccount) + ", "
		if score < 0 {
			query += strconv.Itoa(-1) + ")"
		} else {
			query += strconv.Itoa(1) + ")"
		}
	} else {
		query = "UPDATE vote_comment_to SET vote = "
		x := ""
		if (score == -1 && voteTo.Vote == 1) || (score == 1 && voteTo.Vote == -1) {
			x += strconv.Itoa(0)
		} else if score > 0 {
			x += strconv.Itoa(1)
		} else if score < 0 {
			x += strconv.Itoa(-1)
		} else {
			x += "0"
		}
		query += x + " WHERE id_comment = " + strconv.Itoa(idComment) + " AND id_account = " + strconv.Itoa(idAccount)
	}
	sqlServ.executeQuery(query)
}
