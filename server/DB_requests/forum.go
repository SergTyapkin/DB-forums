package DB_requests

import (
	. "DB-forums/models"
	"github.com/jackc/pgx"
)

func INSERTForum(structureInsert Forum) (Forum, error) {
	var structure Forum
	err := DB.QueryRow(`INSERT INTO Forums(slug, title, author) VALUES ($1, $2, $3) RETURNING *;`,
		structureInsert.Slug, structureInsert.Title, structureInsert.User).
		Scan(&structure.Slug, &structure.Title, &structure.User, &structure.Threads, &structure.Posts)
	return structure, err
}

func INSERTForumToUser(slug string, user User) (User, error) {
	var structure User
	err := DB.QueryRow(`INSERT INTO Forums_to_users(slug, nickname, name, about, email) VALUES ($1, $2, $3, $4, $5) RETURNING *;`,
		slug, user.Nickname, user.Name, user.About, user.Email).
		Scan(&slug, &structure.Nickname, &structure.Name, &structure.Email, &structure.About)
	return structure, err
}

func SELECTForum_slug(slug string) (Forum, error) {
	var structure Forum
	err := DB.QueryRow(`SELECT * FROM Forums WHERE LOWER(slug)=LOWER($1) LIMIT 1;`, slug).
		Scan(&structure.Slug, &structure.Title, &structure.User, &structure.Threads, &structure.Posts)
	return structure, err
}

func INSERTThread(structureInsert Thread) (Thread, error) {
	var structure Thread
	var err error
	if structureInsert.Slug.String != "" {
		err = DB.QueryRow(`INSERT INTO Threads(title, author, created, forum, message, slug) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
			structureInsert.Title, structureInsert.Author, structureInsert.Created, structureInsert.Forum, structureInsert.Message, structureInsert.Slug).
			Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
	} else {
		err = DB.QueryRow(`INSERT INTO Threads(title, author, created, forum, message) VALUES ($1, $2, $3, $4, $5) RETURNING *;`,
			structureInsert.Title, structureInsert.Author, structureInsert.Created, structureInsert.Forum, structureInsert.Message).
			Scan(&structure.Id, &structure.Forum, &structure.Author, &structure.Created, &structure.Message, &structure.Title, &structure.Votes, &structure.Slug)
		if err != nil {
			if structure.Slug.String == "" {
				err = nil
			}
		}
	}
	return structure, err
}

func SELECTForumUsers(slug string, limit int, since string, desc bool) ([]User, error) {
	var rows *pgx.Rows
	var err error
	var users []User

	rows, err = DB.Query(`SELECT nickname, name, about, email
							  FROM forums_to_users
							  WHERE LOWER(slug)=LOWER($1) `+sinceToString("AND nickname", ">", "<", "'", since, "'", desc)+` 
							  ORDER BY nickname `+descToString(desc)+` LIMIT $2;`,
		slug, limit)
	if err != nil {
		return nil, err
	}

	var user User
	for rows.Next() {
		err = rows.Scan(&user.Nickname, &user.Name, &user.About, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, err
}

func SELECTForumThreads(slug string, limit int, since string, desc bool) ([]Thread, error) {
	var rows *pgx.Rows
	var err error
	var threads []Thread

	rows, err = DB.Query(`SELECT *
							  FROM threads
							  WHERE forum=$1 `+sinceToString("AND created", ">=", "<=", "'", since, "'", desc)+` 
							  ORDER BY created `+descToString(desc)+` LIMIT $2;`,
		slug, limit)
	if err != nil {
		return nil, err
	}

	var thread Thread
	for rows.Next() {
		err = rows.Scan(&thread.Id, &thread.Forum, &thread.Author, &thread.Created, &thread.Message, &thread.Title, &thread.Votes, &thread.Slug)
		if err != nil {
			return threads, err
		}
		threads = append(threads, thread)
	}
	defer rows.Close()
	return threads, err
}

func IncrementForumThreads_slug(slug string) (Forum, error) {
	var structure Forum
	err := DB.QueryRow(`UPDATE Forums SET threads = (threads + 1) WHERE LOWER(slug)=LOWER($1) RETURNING *;`, slug).
		Scan(&structure.Slug, &structure.Title, &structure.User, &structure.Threads, &structure.Posts)
	return structure, err
}

func AddForumPosts_slug(slug string, amount int) (Forum, error) {
	var structure Forum
	err := DB.QueryRow(`UPDATE Forums SET posts = (posts + $2) WHERE LOWER(slug)=LOWER($1) RETURNING *;`, slug, amount).
		Scan(&structure.Slug, &structure.Title, &structure.User, &structure.Threads, &structure.Posts)
	return structure, err
}
