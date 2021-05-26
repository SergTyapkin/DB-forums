package DB_requests

import (
	. "DB-forums/models"
	"github.com/jackc/pgx"
)

func SELECTThread_id(id int) (Thread, error) {
	var structure Thread
	err := DB.QueryRow(`SELECT * FROM Threads WHERE id=$1 LIMIT 1;`, id).
		Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
	return structure, err
}

func SELECTThread_slug(slug string) (Thread, error) {
	var structure Thread
	err := DB.QueryRow(`SELECT * FROM Threads WHERE LOWER(slug)=LOWER($1) LIMIT 1;`, slug).
		Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
	return structure, err
}

func UPDATEThread_id(id int, title, message string) (Thread, error) {
	var structure Thread
	err := DB.QueryRow(`UPDATE Threads SET id=id`+
		addIfNotNull("title=", title)+
		addIfNotNull("message=", message)+
		` WHERE id=$1 RETURNING *;`, id).
		Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
	return structure, err
}

func UPDATEThreadVotes_id(id, result int) (Thread, error) {
	var structure Thread
	err := DB.QueryRow(`UPDATE Threads SET votes = (votes + $2) WHERE id=$1 RETURNING *;`, id, result).
		Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
	return structure, err
}

func UPDATEThread_slug(slug string, title, message string) (Thread, error) {
	var structure Thread
	err := DB.QueryRow(`UPDATE Threads SET id=id`+
		addIfNotNull("title=", title)+
		addIfNotNull("message=", message)+
		` WHERE LOWER(slug)=LOWER($1) RETURNING *;`, slug).
		Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
	return structure, err
}

func SELECTThreadPosts_id(id int, limit int, since, sort string, desc bool) ([]Post, error) {
	var rows *pgx.Rows
	var err error

	switch sort {
	case "flat", "":
		rows, err = DB.Query(`SELECT *
							  FROM Posts
							  WHERE thread=$1 `+sinceToString("AND id", ">", "<", "'", since, "'", desc)+` 
							  ORDER BY created `+descToString(desc)+`, id `+descToString(desc)+` LIMIT $2;`,
			id, limit)
	case "tree":
		rows, err = DB.Query(`SELECT *
							  FROM Posts
							  WHERE thread=$1 `+sinceToString("AND paths", ">", "<", "(SELECT paths FROM Posts WHERE id = '", since, "')", desc)+`
							  ORDER BY paths `+descToString(desc)+` LIMIT $2;`,
			id, limit)
	case "parent_tree":
		rows, err = DB.Query(`SELECT *
							  FROM Posts WHERE paths[1] IN
							  (SELECT id FROM posts WHERE thread=$1 AND parent = 0 `+sinceToString("AND paths[1]", ">", "<", "(SELECT paths[1] FROM posts WHERE id = '", since, "')", desc)+` 
							  ORDER BY id `+descToString(desc)+` LIMIT $2)
							  ORDER BY paths[1] `+descToString(desc)+`, paths;`,
			id, limit)
	default:
		return nil, pgx.PgError{Message: "Unknown sort: " + sort}
	}

	if err != nil {
		return nil, err
	}

	var posts []Post
	var post Post
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Author, &post.Created, &post.Forum, &post.Thread, &post.Edited, &post.Message, &post.Parent, &post.Paths)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	defer rows.Close()
	return posts, err
}
