package data

import (
	"context"
	"fmt"
	"time"
)

type Todo struct {
	ID         int       `json:"id,omitempty"`
	UserID     int       `json:"user_id,omitempty"`
	PriorityID int       `json:"priority_id,omitempty"`
	Text       string    `json:"text,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	Priority   Priority  `json:"priority,omitempty"`
}

func (t *Todo) Insert() error {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

	query := "INSERT INTO todos(user_id, priority_id, text, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close() // Ensure the result set is closed after function execution

	result, err := stmt.ExecContext(ctx, t.UserID, t.PriorityID, t.Text, time.Now(), time.Now())
	if err != nil {
		return err
	}

	_, err = result.LastInsertId() // Not interested yet since I reload in the frontend, but I might return the newly inserted record in the future

	return err
}

func (t *Todo) Update() error {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

	query := "UPDATE todos SET priority_id = ?, text = ?, updated_at = ? WHERE id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close() // Ensure the result set is closed after function execution

	_, err = stmt.ExecContext(ctx, t.PriorityID, t.Text, time.Now(), t.ID)

	return err
}

func (t *Todo) GetAll(userID int) ([]Todo, error) {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

	// SQL query with LEFT JOIN on priority table to get complete priority info
	query := `
        SELECT 
            t.id, t.user_id, t.priority_id, t.text, t.created_at, t.updated_at,
            p.id AS priority_id, p.name AS priority_name, p.badge AS priority_badge, p.created_at AS priority_created_at, p.updated_at AS priority_updated_at
        FROM todos t
        LEFT JOIN priorities p ON t.priority_id = p.id
        WHERE t.user_id = ?`

	// Execute the query
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed when function exits

	var todos []Todo

	// Iterate over the query results
	for rows.Next() {
		var todo Todo

		// Scan the current row into the Todo struct, including Priority
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.PriorityID,
			&todo.Text,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Priority.ID,
			&todo.Priority.Name,
			&todo.Priority.Badge,
			&todo.Priority.CreatedAt,
			&todo.Priority.UpdatedAt,
		)
		if err != nil {
			return nil, err // Return an error if scanning fails
		}

		// Append the todo to the list of todos
		todos = append(todos, todo)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *Todo) Delete(ID, userID int) error {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

	// Define the query
	query := "DELETE FROM todos WHERE id = ? AND user_id = ?"
	// Prepare the statement
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close() // Ensure the result set is closed after function execution

	// Execute the statement
	result, err := stmt.ExecContext(ctx, ID, userID)
	if err != nil {
		return err
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the deletion was not successful
	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}

	// If rows were affected, the deletion was successful
	return err
}
