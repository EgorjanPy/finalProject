package tools

import (
	"errors"
	"net/http"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return "", errors.New("userID not found")
	}
	return userID, nil
}
