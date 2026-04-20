package accountgreeter

import "fmt"

// formatName concatenates first and last name with a space.
// This function is unexported — it can only be used within this package.
// (Lowercase first letter = package-private in Go.)
//
// TODO: implement this function.
func formatName(first, last string) string {
	return first + " " + last
}

// Greet returns a personalised welcome message.
// This function is exported — it can be called from any package.
// (Uppercase first letter = exported/public in Go.)
//
// Expected format: "Hello, <First Last>! Your account <accountID> is ready."
//
// TODO: implement using formatName and fmt.Sprintf.
func Greet(accountID, first, last string) string {
	return fmt.Sprintf("Hello, %s! Your account %s is ready.", formatName(first, last), accountID)
}
