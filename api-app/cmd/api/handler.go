package main

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
type envelope map[string]any


type ValidationError struct {
	Errors map[string]string
}

func (v *ValidationError) Error() string {
	return "There was an issue with the validation process."
}

