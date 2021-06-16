package DB_requests

import (
	. "DB-forums/models"
	"fmt"
	"github.com/jackc/pgx"
	"strings"
)

func SELECTPost_id(id int) (Post, error) {
	var structure Post
	err := DB.QueryRow(`SELECT * FROM Posts WHERE id=$1 LIMIT 1;`, id).
		Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
	return structure, err
}

func UPDATEPost_id(id int, message string) (Post, error) {
	var structure Post
	err := DB.QueryRow(`UPDATE Posts SET edited = TRUE`+addIfNotNull(`message = `, message)+` WHERE id=$1 RETURNING *;`, id).
		Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
	return structure, err
}

func INSERTPosts(postsInsert []Post, thread int, forum string) ([]Post, error) {
	var posts []Post
	var err error

	queryInsertUsers := `INSERT INTO Forums_to_users(slug, nickname, name, about, email)
						 SELECT '` + forum + `', nickname, name, about, email
						 FROM Users WHERE nickname IN (`
	queryInsert := `INSERT INTO Posts(author, created, forum, thread, message, parent) VALUES `
	var valuesInsert []interface{}
	authorsSet := make(map[string]bool)
	base := 0
	for _, post := range postsInsert {
		if post.Created.IsZero() {
			queryInsert += fmt.Sprintf(`($%d, NOW(), $%d, $%d, $%d, $%d), `, base+1, base+2, base+3, base+4, base+5)
			valuesInsert = append(valuesInsert, post.Author, forum, thread, post.Message, post.Parent)
			base += 5
		} else {
			queryInsert += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d), `, base+1, base+2, base+3, base+4, base+5, base+6)
			valuesInsert = append(valuesInsert, post.Author, post.Created, forum, thread, post.Message, post.Parent)
			base += 6
		}

		if _, ex := authorsSet[post.Author]; !ex { // если такого автора ещё нет
			queryInsertUsers += `'` + post.Author + `', `
			authorsSet[post.Author] = true
		}
	}
	queryInsert = strings.TrimSuffix(queryInsert, `, `) + ` RETURNING *;`
	queryInsertUsers = strings.TrimSuffix(queryInsertUsers, `, `) + `) ON CONFLICT DO NOTHING;`

	// insert posts
	rows, err := DB.Query(queryInsert, valuesInsert...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var structure Post
		err = rows.Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
		if err != nil {
			return nil, err
		}
		posts = append(posts, structure)
	}

	if pgErr, ok := rows.Err().(pgx.PgError); ok {
		if pgErr.Code != "42601" { // Кривой запрос (когда нет записей)
			return nil, rows.Err()
		}
	}

	// Добавляем в forums_to_users
	_, err = DB.Exec(queryInsertUsers)
	if pgErr, ok := err.(pgx.PgError); ok {
		if pgErr.Code == "42601" { // Кривой запрос (когда нет записей)
			return posts, nil
		}
	}

	return posts, err
}
