package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db?_foreign_keys=1")
	if err != nil {
		panic("Could not connect to the database")
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(50)
	DB.SetConnMaxLifetime(time.Minute * 5)

	createTable()
}

func createTable() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		fmt.Println(err)
		panic("Could not create users table.")
	}

	createPriorityTable := `
	CREATE TABLE IF NOT EXISTS priorities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		badge TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = DB.Exec(createPriorityTable)
	if err != nil {
		fmt.Println(err)
		panic("Could not create priorities table.")
	}

	// Insert default priorities if they don’t exist
	insertPriorities := `
	INSERT INTO priorities (name, badge, created_at, updated_at)
	SELECT 'Low', 'is-info', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	WHERE NOT EXISTS (SELECT 1 FROM priorities WHERE name = 'Low')
	UNION ALL
	SELECT 'Medium', 'is-warning', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	WHERE NOT EXISTS (SELECT 1 FROM priorities WHERE name = 'Medium')
	UNION ALL
	SELECT 'High', 'is-danger', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	WHERE NOT EXISTS (SELECT 1 FROM priorities WHERE name = 'High')
	`

	_, err = DB.Exec(insertPriorities)
	if err != nil {
		fmt.Println(err)
		panic("Could not insert default priorities.")
	}

	createTodoTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		priority_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (priority_id) REFERENCES priorities(id) ON DELETE SET NULL
	)`

	_, err = DB.Exec(createTodoTable)
	if err != nil {
		fmt.Println(err)
		panic("Could not create todos table.")
	}
}
