package DB_requests

import (
	. "DB-forums/models"
)

func INSERTUser(structureInsert User) (User, error) {
	var structure User
	err := DB.QueryRow(`INSERT INTO Users(Nickname, Name, Email, About) VALUES ($1, $2, $3, $4) RETURNING *;`,
		structureInsert.Nickname, structureInsert.Name, structureInsert.Email, structureInsert.About).
		Scan(&structure.Nickname, &structure.Name, &structure.Email, &structure.About)
	return structure, err
}

func SELECTUser_nickname(nickname string) (User, error) {
	var structure User
	row := DB.QueryRow(`SELECT * FROM Users WHERE nickname=$1 LIMIT 1;`, nickname)
	err := row.Scan(&structure.Nickname, &structure.Name, &structure.Email, &structure.About)
	return structure, err
}

func SELECTUser_email(email string) (User, error) {
	var structure User
	row := DB.QueryRow(`SELECT * FROM Users WHERE LOWER(email)=LOWER($1) LIMIT 1;`, email)
	err := row.Scan(&structure.Nickname, &structure.Name, &structure.Email, &structure.About)
	return structure, err
}

func UPDATEUser_nickname(nickname, newNickname, name, about, email string) (User, error) {
	var structure User
	err := DB.QueryRow(`UPDATE Users SET nickname=$1`+
		addIfNotNull("name=", name)+
		addIfNotNull("email=", email)+
		addIfNotNull("about=", about)+
		` WHERE nickname=$1 RETURNING *;`, newNickname).
		Scan(&structure.Nickname, &structure.Name, &structure.Email, &structure.About)
	return structure, err
}
