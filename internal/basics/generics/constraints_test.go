package generics

import (
	"net/http"
	"testing"
)

func TestGenericsToString(t *testing.T) {
	m1 := Magic{"Protection Spell", []string{"Expecto", "Patronum"}}
	m2 := Magic{"meme", []string{"Expensive", "Petroleum"}}
	err := CustomError{http.StatusNotImplemented, "Panic! Your code is on fire!"}

	ToString(&m1)
	ToString(&m2)
	ToString(&err)
}
