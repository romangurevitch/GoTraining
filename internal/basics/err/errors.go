package err

import (
	"errors"
	"fmt"
)

// --- Sentinel Errors ---

var ErrNotFound = errors.New("item not found")
var ErrPermission = errors.New("permission denied")

// --- Custom Error Type ---

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Reason)
}

// --- Example functions that return errors ---

// LookupUser returns ErrNotFound if the ID is unknown.
func LookupUser(id string) (string, error) {
	users := map[string]string{
		"1": "Alice",
		"2": "Bob",
	}
	name, ok := users[id]
	if !ok {
		return "", fmt.Errorf("lookup user %s: %w", id, ErrNotFound)
	}
	return name, nil
}

// GreetUser looks up a user and returns a greeting.
// Shows the typical if err != nil pattern for checking and wrapping errors.
func GreetUser(id string) (string, error) {
	name, err := LookupUser(id)
	if err != nil {
		return "", fmt.Errorf("greet: %w", err)
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}
