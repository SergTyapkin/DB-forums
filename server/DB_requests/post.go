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

func INSERTPosts(structuresInsert []Post, thread int, forum string) ([]Post, error) {
	var structures []Post
	var err error
	query := `INSERT INTO Posts(author, created, forum, thread, message, parent) VALUES `
	var values []interface{}
	for i, structureInsert := range structuresInsert {
		base := i * 6
		query += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d), `, base+1, base+2, base+3, base+4, base+5, base+6)
		values = append(values, structureInsert.Author, structureInsert.Created, forum, thread, structureInsert.Message, structureInsert.Parent)
	}
	query = strings.TrimSuffix(query, `, `) + ` RETURNING *;`

	rows, err := DB.Query(query, values...)
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
		structures = append(structures, structure)
	}

	if pgErr, ok := rows.Err().(pgx.PgError); ok {
		if pgErr.Code != "42601" { // Кривой запрос (когда нет записей)
			return nil, rows.Err()
		}
	}

	return structures, err
}
