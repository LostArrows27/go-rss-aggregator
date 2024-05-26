package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey is a function that extracts the API key from the headers
// Example:
// Authorization: ApiKey {API_KEY}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no API key found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("invalid API key format")
	}

	return vals[1], nil
}
