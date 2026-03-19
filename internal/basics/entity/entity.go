package entity

import (
	"fmt"
	"strings"
)

// User is the Interface representing the behaviors our entity must have.
type User interface {
	GetID() string
	GetName() string
	GetEmail() string
	IsAdmin() bool
	String() string
}

// user is the concrete struct (unexported for encapsulation).
// Users of this package should interact with the User interface.
type user struct {
	ID    string `json:"id"`    // Exported field
	Name  string `json:"name"`  // Exported field
	Email string `json:"email"` // Exported field
	role  string // Unexported field (internal behavior)
}

// New is a factory function that creates a new User entity.
func New(id, name, email string) User {
	return &user{
		ID:    id,
		Name:  name,
		Email: strings.ToLower(email),
		role:  "customer", // Default role
	}
}

// --- Methods implementing the User interface ---

func (u *user) GetID() string {
	return u.ID
}

func (u *user) GetName() string {
	return u.Name
}

func (u *user) GetEmail() string {
	return u.Email
}

func (u *user) IsAdmin() bool {
	return u.role == "admin"
}

// String provides a string representation of the entity (implements fmt.Stringer)
func (u *user) String() string {
	return fmt.Sprintf("User[%s]: %s <%s>", u.ID, u.Name, u.Email)
}

// Promotion (to demonstrate that we can change state via unexported fields internally)
func (u *user) PromoteToAdmin() {
	u.role = "admin"
}
