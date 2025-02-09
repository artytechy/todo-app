package main

import (
	"errors"
	"net/http"
	"task-app/db/data"
	"task-app/utils"
)

type ValidationErrors struct {
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirm_password,omitempty"`
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.User.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "OK",
		Data:    envelope{"users": users},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = validateRegisterUserInputs(app, &user)
	if err != nil {
		if validationErr, ok := err.(*ValidationError); ok {
			// Map validation errors as needed
			app.errorJSONWithData(w, err, envelope{"errors": validationErr.Errors})
			return
		}
	}

	_, err = app.models.User.Insert(user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Welcome! Your registration was successful.",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *application) LoginUser(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var creds credentials
	var payload jsonResponse
	var userValidationErrors = map[string]string{}

	err := app.readJSON(w, r, &creds)
	if err != nil {
		payload.Error = true
		payload.Message = "Oops! Something went wrong. Please try again later."
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	// Check for empty fields and add error messages to map
	checkEmptyField(&userValidationErrors, "email", creds.Email)
	checkEmptyField(&userValidationErrors, "password", creds.Password)
	if len(userValidationErrors) > 0 {
		payload.Error = true
		payload.Message = "There was an issue with the validation process."
		payload.Data = envelope{"errors": userValidationErrors}
		_ = app.writeJSON(w, http.StatusBadRequest, payload)

		return
	}

	
	user, err := app.models.User.GetByEmail(creds.Email)
	if err != nil {
		payload.Error = true
		payload.Message = "Authentication failed."
		_ = app.writeJSON(w, http.StatusBadRequest, payload)

		return
	}

	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		payload.Error = true
		payload.Message = "Authentication failed."
		_ = app.writeJSON(w, http.StatusBadRequest, payload)

		return
	}

	token, err := utils.GenerateToken(user.Email, int64(user.ID))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload = jsonResponse{
		Error:   false,
		Message: "Welcome! It's great to see you again!",
		Data: envelope{"user": user, "token": token},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) LogoutUser(w http.ResponseWriter, r *http.Request) {
		// Extract the userID from the request context (set by Authenticate middleware)
		_, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSONWithData(w, errors.New("invalid json"), requestPayload)
		
		return
	}

	utils.InvalidateToken(requestPayload.Token)

	payload := jsonResponse{
		Error:   false,
		Message: "You've been logged out successfully!",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func validateRegisterUserInputs(app *application, user *data.User) error {
	var userValidationErrors = map[string]string{}

	// Check for empty fields and add error messages to map
	checkEmptyField(&userValidationErrors, "name", user.Name)
	checkEmptyField(&userValidationErrors, "email", user.Email)
	checkEmptyField(&userValidationErrors, "password", user.Password)
	checkEmptyField(&userValidationErrors, "confirm_password", user.ConfirmPassword)

	// Check if email already exists
	if user.Email != "" {
		emailExists, err := app.models.User.EmailExists(user.Email)
		if err != nil {
			userValidationErrors["email"] = "An error occurred while checking email availability. Please try again later."
		} else if emailExists {
			userValidationErrors["email"] = "It looks like this email is already in use. Try another one."
		}
	}

	// Ensure password and confirm password match
	if user.Password != "" && user.ConfirmPassword != "" {
		if user.Password != user.ConfirmPassword {
			userValidationErrors["confirm_password"] = "The passwords you entered don’t match. Please try again."
		}
	}


	// If there are validation errors, return a ValidationError with the error map
	if len(userValidationErrors) > 0 {
		return &ValidationError{Errors: userValidationErrors}
	}

	// No validation errors
	return nil
}
