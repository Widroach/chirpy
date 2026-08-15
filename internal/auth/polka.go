package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "ApiKey")
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("error no \"Authorization: ApiKey <key>\" header is missing or invalid")
	}
	return token, nil
}
