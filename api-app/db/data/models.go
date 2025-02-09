package data

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
		Priority: Priority{},
		Todo: Todo{},
	}
}

type Models struct {
	User User
	Priority Priority
	Todo Todo
}
