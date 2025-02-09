package main

import (
	"net/http"
	"task-app/db/data"
)

func (app *application) SaveTodo(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userID := r.Context().Value(userIDKey).(int64)
	if userID == 0 {
		// handle missing userID case
		app.writeJSON(w, http.StatusUnauthorized, jsonResponse{
			Error:   true,
			Message: "Unauthorized.",
		})
		return
	}

	var requestPayload struct {
		ID         int    `json:"id"`
		PriorityID int    `json:"priority_id"`
		Text       string `json:"text"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	todo := data.Todo{
		ID:         requestPayload.ID,
		UserID:     int(userID),
		PriorityID: requestPayload.PriorityID,
		Text:       requestPayload.Text,
	}

	err = validateTodoInputs(&todo)
	if err != nil {
		if validationErr, ok := err.(*ValidationError); ok {
			// Map validation errors as needed
			app.errorJSONWithData(w, err, envelope{"errors": validationErr.Errors})
			return
		}
	}

	if todo.ID == 0 {
		err = todo.Insert()
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = todo.Update()
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Todo has been successfully saved.",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *application) AllTodos(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userID := r.Context().Value(userIDKey).(int64)
	if userID == 0 {
		// handle missing userID case
		app.writeJSON(w, http.StatusUnauthorized, jsonResponse{
			Error:   true,
			Message: "Unauthorized.",
		})

		return
	}

	todos, err := app.models.Todo.GetAll(int(userID))
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "OK",
		Data:    envelope{"todos": todos},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userID := r.Context().Value(userIDKey).(int64)
	if userID == 0 {
		// handle missing userID case
		app.writeJSON(w, http.StatusUnauthorized, jsonResponse{
			Error:   true,
			Message: "Unauthorized.",
		})

		return
	}
	var requestPayload struct {
		ID int `json:"id"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.Todo.Delete(requestPayload.ID, int(userID))
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, jsonResponse{
			Error:   true,
			Message: "Whoops! Something went wrong. Please try again later..",
		})
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Todo has been successfully deleted.",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func validateTodoInputs(todo *data.Todo) error {
	var todoValidationErrors = map[string]string{}

	// Check for empty fields and add error messages to map
	checkEmptyField(&todoValidationErrors, "priority_id", todo.PriorityID)
	checkEmptyField(&todoValidationErrors, "text", todo.Text)

	// If there are validation errors, return a ValidationError with the error map
	if len(todoValidationErrors) > 0 {
		return &ValidationError{Errors: todoValidationErrors}
	}

	// No validation errors
	return nil
}
