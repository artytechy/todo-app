package main

import (
	"context"
	"net/http"
	"strings"
	"task-app/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

func (app *application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			payload := jsonResponse{
				Error:   true,
				Message: "Authorization header missing.",
			}
			app.writeJSON(w, http.StatusUnauthorized, payload)
			return
		}

		// Ensure header contains exactly two parts: "Bearer <token>"
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			payload := jsonResponse{
				Error:   true,
				Message: "Invalid authorization format. Expected 'Bearer <token>'.",
			}
			app.writeJSON(w, http.StatusUnauthorized, payload)
			return
		}

		token := headerParts[1]
		userID, err := utils.VerifyToken(token)
		if err != nil {
			payload := jsonResponse{
				Error:   true,
				Message: "Invalid or expired token.",
				Data: headerParts,
			}
			app.writeJSON(w, http.StatusUnauthorized, payload)
			return
		}

		// Store user ID in request context for later use
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
