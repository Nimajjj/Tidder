package mySQL

import (
	"strconv"
	"time"

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
	if conditions == " ORDER BY creation_date DESC" {
		query += " ORDER BY creation_date DESC" // LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOL FLEMME
	} else if conditions != "" {
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

func (sqlServ SqlServer) CreatePost(title string, media_url string, content string, nsfw_input bool, id_subject string, id_author int) {
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

	if content == "" && media_url == "" {
		Util.Warning("Creating sub failed : content and media cannot be both empty.")
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

func (sqlServ SqlServer) Vote(id_post int, score int, id_account int) {
	Util.Log("Vote received: " + strconv.Itoa(id_post) + " " + strconv.Itoa(score))
	query := "UPDATE tidder.posts SET "
	if score < 0 {
		query += "downvotes = downvotes + " + strconv.Itoa(score*(-1))
	} else {
		query += "upvotes = upvotes + " + strconv.Itoa(score)
	}
	query += " WHERE id_post = " + strconv.Itoa(id_post)
	sqlServ.executeQuery(query)

	type VoteTo struct {
		IdPost    int
		IdAccount int
		Vote      int
		Shit      int
	}

	var voteTo VoteTo
	query = "SELECT * FROM vote_to WHERE id_account = " + strconv.Itoa(id_account) + " AND id_post = " + strconv.Itoa(id_post)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	res := 0
	for rows.Next() {
		res += 1
		var id_post int
		var id_account int
		var vote int
		var shit int
		if err2 := rows.Scan(
			&shit,
			&id_post,
			&id_account,
			&vote,
		); err2 != nil {
			Util.Error(err2)
		}
		voteTo = VoteTo{id_post, id_account, vote, shit}
	}
	if res == 0 {
		query = "INSERT INTO vote_to (id_post, id_account, vote) VALUES (" + strconv.Itoa(id_post) + ", " + strconv.Itoa(id_account) + ", "
		if score < 0 {
			query += strconv.Itoa(-1) + ")"
		} else {
			query += strconv.Itoa(1) + ")"
		}
	} else {
		query = "UPDATE vote_to SET vote = "
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
		query += x + " WHERE id_post = " + strconv.Itoa(id_post) + " AND id_account = " + strconv.Itoa(id_account)
	}
	sqlServ.executeQuery(query)
}

func (sqlServ SqlServer) MakeDisplayablePost(post Posts, account_id int) DisplayablePost {
	vote := 0
	if account_id != -1 {
		query := "SELECT * FROM vote_to WHERE id_account = " + strconv.Itoa(sqlServ.GetAccountById(account_id).Id) + " AND id_post = " + strconv.Itoa(post.Id)
		Util.Query(query)
		rows, err := sqlServ.db.Query(query)
		if err != nil {
			Util.Error(err)
		}
		for rows.Next() {
			var id_post int
			var id_account int
			var shit int
			if err2 := rows.Scan(
				&shit,
				&id_post,
				&id_account,
				&vote,
			); err2 != nil {
				Util.Error(err2)
			}
		}

	}

	sub := sqlServ.GetSubs("id_subject = " + strconv.Itoa(post.IdSubject))[0]
	return DisplayablePost{
		post.Id,
		post.Title,
		post.MediaUrl,
		post.Content,
		post.CreationDate,
		post.Nsfw,
		post.Redacted,
		post.Pinned,
		post.Score,
		post.NumberOfComments,
		vote,
		sqlServ.GetAccountById(post.IdAuthor).Name,
		sub.Name,
		sub.ProfilePicture,
	}
}

func (sqlServ SqlServer) GenerateFeed(user int) []DisplayablePost {
	query := ""
	if user != -1 {
		subtidders := sqlServ.GetSubtiddersSubscribed(user)
		for i, subtidder := range subtidders {
			query += "id_subject = " + strconv.Itoa(subtidder.Id) + " "
			if i < len(subtidders)-1 {
				query += "OR "
			}
		}
	}

	query += " ORDER BY creation_date DESC"

	res := []DisplayablePost{}
	for _, post := range sqlServ.GetPosts(query) {
		res = append(res, sqlServ.MakeDisplayablePost(post, user))
	}

	return res
}

func (sqlServ SqlServer) GenerateSubTidderFeed(user int, subtidder int) []DisplayablePost {
	query := " id_subject = " + strconv.Itoa(subtidder) + " ORDER BY creation_date DESC"

	res := []DisplayablePost{}
	for _, post := range sqlServ.GetPosts(query) {
		res = append(res, sqlServ.MakeDisplayablePost(post, user))
	}

	return res
}
