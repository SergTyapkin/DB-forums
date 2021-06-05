package DB_requests

import (
	. "DB-forums/models"
)

func INSERTVote(structureInsert Vote) (Vote, error) {
	var structure Vote
	boolResult := true
	if structureInsert.Result == -1 {
		boolResult = false
	}
	err := DB.QueryRow(`INSERT INTO Votes(nickname, result, thread) VALUES ($1, $2, $3) RETURNING *;`,
		structureInsert.Nickname, boolResult, structureInsert.Thread).
		Scan(&structure.Nickname, &boolResult, &structure.Thread)
	structure.Result = 1
	if boolResult == false {
		structure.Result = 0
	}
	return structure, err
}

func SELECTVote_nickname_thread(nickname string, thread_id int) (Vote, error) {
	var structure Vote
	var boolResult bool
	err := DB.QueryRow(`SELECT * FROM Votes WHERE nickname = $1 AND thread=$2;`, nickname, thread_id).
		Scan(&structure.Nickname, &boolResult, &structure.Thread)
	structure.Result = 1
	if boolResult == false {
		structure.Result = -1
	}
	return structure, err
}

func UPDATEVote_nickname_thread(nickname string, thread_id, result int) (Vote, error) {
	var structure Vote
	boolResult := true
	if result == -1 {
		boolResult = false
	}
	err := DB.QueryRow(`UPDATE Votes SET result = $3 WHERE nickname = $1 AND thread=$2 RETURNING *;`, nickname, thread_id, boolResult).
		Scan(&structure.Nickname, &boolResult, &structure.Thread)
	structure.Result = 1
	if boolResult == false {
		structure.Result = -1
	}
	return structure, err
}
