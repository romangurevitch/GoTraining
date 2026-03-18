package err

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myError struct {
	myMessage string
}

func (e myError) Error() string {
	return e.myMessage
}

func TestIs(t *testing.T) {
	specificErr := errors.New("cause 1")
	wrappedErr := fmt.Errorf("something happened because: %w", specificErr)

	sameTextErr := errors.New("cause 1")

	// will errors.Is(wrappedErr, specificErr) return true or false?
	// will errors.Is(wrappedErr, sameTextErr) return true or false?
	assert.True(t, errors.Is(wrappedErr, specificErr))
	assert.False(t, errors.Is(wrappedErr, sameTextErr))
	fmt.Println("err:", wrappedErr)
}

func TestAs(t *testing.T) {
	err := myError{myMessage: "special error"}

	printMessage(err)
}

func printMessage(err error) {

	// will this worK?
	//fmt.Println(err.myMessage)

	var mySpecialError myError
	if errors.As(err, &mySpecialError) {
		fmt.Println("I can access my special error message:", mySpecialError.myMessage)
	}
}

func TestWrappedAs(t *testing.T) {
	specificErr := myError{myMessage: "special error"}
	wrappedErr := fmt.Errorf("something happened because: %w", specificErr)

	var mySpecialError myError
	if errors.As(wrappedErr, &mySpecialError) {
		fmt.Println("I can access my special error message:", mySpecialError.myMessage)
	}
}
