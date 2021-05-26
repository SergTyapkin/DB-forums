package DB_requests

import (
	. "DB-forums/models"
	"github.com/go-openapi/swag"
)

func SELECTPost_id(id int) (Post, error) {
	var structure Post
	err := DB.QueryRow(`SELECT * FROM Posts WHERE id=$1 LIMIT 1;`, id).
		Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
	return structure, err
}

func UPDATEPost_id(id int, message string) (Post, error) {
	var structure Post
	err := DB.QueryRow(`UPDATE Posts SET edited = true`+addIfNotNull("message=", message)+` WHERE id=$1 RETURNING *;`, id).
		Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
	return structure, err
}

func INSERTPost(structureInsert Post) (Post, error) {
	var structure Post
	var err error
	if !swag.IsZero(structureInsert.Created) {

		err = DB.QueryRow(`INSERT INTO Posts(author, created, forum, thread, message, parent) VALUES ($1, $2, $3, $4, $5, $6)
							   RETURNING *;`,
			structureInsert.Author, structureInsert.Created, structureInsert.Forum, structureInsert.Thread, structureInsert.Message, structureInsert.Parent).
			Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
	} else {
		err = DB.QueryRow(`INSERT INTO Posts(author, forum, thread, message, parent) VALUES ($1, $2, $3, $4, $5)
							   RETURNING *;`,
			structureInsert.Author, structureInsert.Forum, structureInsert.Thread, structureInsert.Message, structureInsert.Parent).
			Scan(&structure.Id, &structure.Author, &structure.Created, &structure.Forum, &structure.Thread, &structure.Edited, &structure.Message, &structure.Parent, &structure.Paths)
	}
	if err != nil {
		return structure, err
	}
	return structure, err
}
