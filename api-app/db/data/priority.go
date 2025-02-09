package data

import (
	"context"
	"time"
)

type Priority struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Badge     string    `json:"badge,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (p *Priority) GetAll() ([]Priority, error) {
	// Create a new context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel() // Ensure the context is canceled when the function exits

		// Define the SQL query to retrieve all priorities from the database
	query := "SELECT * FROM priorities"

	// Execute the query using the context to ensure it respects the timeout
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Return an error if the query fails
	}
	defer rows.Close() // Ensure the result set is closed after function execution

	// Slice to store the list of priorities retrieved from the database
	var priorities []Priority

	// Iterate over the query results
	for rows.Next() {
		var priority Priority

		// Scan the current row into the priority struct fields
		err := rows.Scan(&priority.ID, &priority.Name, &priority.Badge,  &priority.CreatedAt, &priority.UpdatedAt)
		if err != nil {
			return nil, err // Return an error if scanning fails
		}

		// Append the priority to the list of priorities
		priorities = append(priorities, priority)
	}

	// Check if any error occurred during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of users and no error
	return priorities, nil
}