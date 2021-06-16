package models

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"

	"time"
)

var DB *pgx.ConnPool

type User struct {
	Nickname string `json:"nickname"`
	Name     string `json:"fullname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}

type Forum struct {
	Slug    string `json:"slug"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Threads int    `json:"threads"`
	Posts   int    `json:"posts"`
}

type Thread struct {
	Id      int       `json:"id"`
	Forum   string    `json:"forum"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Message string    `json:"message"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes"`
	Slug    NullStr   `json:"slug"`
}

type Post struct {
	Id       int64            `json:"id"`
	Author   string           `json:"author"`
	Created  time.Time        `json:"created"`
	Forum    string           `json:"forum"`
	Thread   int              `json:"thread"`
	Edited   bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   int              `json:"parent"`
	Paths    pgtype.Int8Array `json:"-"`
	IntPaths []int            `json:"-"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Result   int    `json:"voice"`
	Thread   int    `json:"thread"`
}

//---------------------------

type Stats struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}

type JsonError struct {
	Message string `json:"message"`
}

type PostUpdate struct {
	Message string `json:"message"`
}
