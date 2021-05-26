package DB_requests

import (
	. "DB-forums/models"
	"io/ioutil"
	"os"
)

func InitDB_all() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "Can't get filepath", err
	}
	data, errRead := ioutil.ReadFile(pwd + "/sql/initial.sql")
	if errRead != nil {
		return "Can't read file: " + pwd + "/sql/initial.sql", err
	}
	_, err = DB.Exec(string(data))
	if err != nil {
		return "Invalid database request: " + string(data), err
	}
	return "Database initialized", err
}

func DropDB_all() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "Can't get filepath", err
	}
	data, errRead := ioutil.ReadFile(pwd + "/sql/drop.sql")
	if errRead != nil {
		return "Can't read file: " + pwd + "/sql/drop.sql", err
	}
	_, err = DB.Exec(string(data))
	if err != nil {
		return "Invalid database request: " + string(data), err
	}
	return "Database cleared", err
}

func ClearDB_all() (string, error) {
	_, err := DB.Exec(`TRUNCATE Forums, Posts, Threads, Users, Votes, Forums_to_users;`)
	if err != nil {
		return "Invalid database request", err
	}
	return "Database cleared", err
}

func GetStats() (Stats, error) {
	var stats Stats
	err := DB.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&stats.User)
	if err != nil {
		stats.User = 0
	}
	err = DB.QueryRow(`SELECT COUNT(*) FROM forums;`).Scan(&stats.Forum)
	if err != nil {
		stats.Forum = 0
	}
	err = DB.QueryRow(`SELECT COUNT(*) FROM threads;`).Scan(&stats.Thread)
	if err != nil {
		stats.Thread = 0
	}
	err = DB.QueryRow(`SELECT COUNT(*) FROM posts;`).Scan(&stats.Post)
	if err != nil {
		stats.Post = 0
	}
	return stats, err
}
