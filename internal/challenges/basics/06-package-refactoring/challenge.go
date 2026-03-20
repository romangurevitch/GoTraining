// Package util is a god-package that needs to be refactored.
// See README.md for the refactoring requirements.
package util

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var defaultLocation *time.Location

// init silently sets a global timezone.
// Requirement: Replace this with an explicit Setup(tz string) error function.
func init() {
	defaultLocation, _ = time.LoadLocation("UTC")
}

// ParseDate parses RFC3339 strings.
// Requirement: Move to a dedicated 'dates' package.
func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation(time.RFC3339, s, defaultLocation)
}

// Sanitize trims whitespace and lowercases.
// Requirement: Move to a dedicated 'sanitize' package.
func Sanitize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// BuildAuthHeader returns "Bearer <token>".
// Requirement: Move to a dedicated 'auth' or 'headers' package.
// BUG: This currently uses the wrong prefix! Fix it during refactor.
func BuildAuthHeader(token string) string {
	return fmt.Sprintf("%s %s", http.MethodGet, token)
}
