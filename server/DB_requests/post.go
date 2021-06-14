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
	/*
		-- Check thread of parent post
		IF (NEW.parent <> 0) THEN
		SELECT thread FROM Posts WHERE id = NEW.parent INTO parent_thread;
		IF (NOT FOUND) OR parent_thread <> NEW.thread THEN
		RAISE EXCEPTION 'Parent post in another thread' USING ERRCODE = '00228';
		END IF;
		END IF;
	*/
	querySelect := `SELECT thread FROM Posts WHERE id IN (`
	queryInsert := `INSERT INTO Posts(author, created, forum, thread, message, parent) VALUES `
	var valuesInsert, valuesSelect []interface{}
	isNeedToCheckParents := false
	selectIndex := 0
	selectParentsSet := make(map[int]bool)
	for i, post := range postsInsert {
		base := i * 6
		queryInsert += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d), `, base+1, base+2, base+3, base+4, base+5, base+6)
		valuesInsert = append(valuesInsert, post.Author, post.Created, forum, thread, post.Message, post.Parent)
		if post.Parent != 0 {
			isNeedToCheckParents = true
			if _, ex := selectParentsSet[post.Parent]; !ex { // если такого родителя ещё нет
				selectParentsSet[post.Parent] = true // добавляем в множество
				selectIndex += 1
				querySelect += fmt.Sprintf(`$%d, `, selectIndex)
				valuesSelect = append(valuesSelect, post.Parent)
			}
		}
	}
	queryInsert = strings.TrimSuffix(queryInsert, `, `) + ` RETURNING *;`

	// check post parents
	if isNeedToCheckParents {
		querySelect = strings.TrimSuffix(querySelect, `, `) + `);`

		rows, err := DB.Query(querySelect, valuesSelect...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		rowsCount := 0
		for rows.Next() {
			var structure PostForCheck
			err = rows.Scan(&structure.Thread)
			if err != nil {
				return nil, err
			}
			rowsCount += 1
			if thread != structure.Thread {
				return nil, pgx.ErrDeadConn
			}
		}
		//println("Count: ", rowsCount, " vs ", selectIndex)
		if rowsCount != selectIndex {
			return nil, pgx.ErrDeadConn
		}
	}

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

	return posts, err
}
