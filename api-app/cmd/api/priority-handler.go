package main

import "net/http"

func (app *application) AllPriorities(w http.ResponseWriter, r *http.Request) {
	priorities, err := app.models.Priority.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "OK",
		Data:    envelope{"priorities": priorities},
	}

	app.writeJSON(w, http.StatusOK, payload)
}