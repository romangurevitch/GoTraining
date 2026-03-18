package main

import (
	"fmt"
	"os"

	greeting "github.com/romangurevitch/go-training/internal/hello"
)

func main() {
	// Extracting a command-line argument if provided by the user
	// os.Args is the program name; os.Args is the first argument
	targetName := ""

	if len(os.Args) > 1 {
		targetName = os.Args[1]
	}
	// Invoking the decoupled domain logic
	message := greeting.Generate(targetName)

	// Executing the side-effect (I/O to the terminal)
	fmt.Println(message)
}
