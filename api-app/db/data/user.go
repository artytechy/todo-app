package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int       `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Email           string    `json:"email,omitempty"`
	Password        string    `json:"password,omitempty"`
	ConfirmPassword string    `json:"confirm_password,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

func (u *User) GetAll() ([]User, error) {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

	// Define the SQL query to retrieve all users from the database
	query := "SELECT * FROM users"

	// Execute the query using the context to ensure it respects the timeout
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Return an error if the query fails
	}
	defer rows.Close() // Ensure the result set is closed after function execution

	// Slice to store the list of users retrieved from the database
	var users []User

	// Iterate over the query results
	for rows.Next() {
		var user User

		// Scan the current row into the user struct fields
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err // Return an error if scanning fails
		}

		// Append the user to the list of users
		users = append(users, user)
	}

	// Check if any error occurred during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of users and no error
	return users, nil
}

func (u *User) Insert(user User) (int, error) {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO users(name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	defer stmt.Close() // Ensure the result set is closed after function execution

	result, err := stmt.Exec(user.Name, user.Email, hashedPassword, time.Now(), time.Now())
	if err != nil {
		return 0, nil
	}

	userID, err := result.LastInsertId()
	u.ID = int(userID)

	return int(userID), err
}

func (u *User) EmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel() // Ensure the context is canceled when the function exits

    // Query to check if the email exists
    query := "SELECT email FROM users WHERE email = ? LIMIT 1"
    row := db.QueryRowContext(ctx, query, email)

    var retrievedEmail string
    if err := row.Scan(&retrievedEmail); err != nil {
        if err == sql.ErrNoRows {
            // No rows found, so the email doesn't exist
            return false, nil
        }
        // Return the error if there's an issue with the query or scanning
        return false, err
    }

    // If retrievedEmail is not empty, the email exists
    return retrievedEmail != "", nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel() // Ensure the context is canceled when the function exits

	// Query to get user by email
    query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ? LIMIT 1"

	var user User
    row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID, 
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) PasswordMatches(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
