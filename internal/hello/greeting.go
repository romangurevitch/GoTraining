package hello

const defaultGreetingPrefix = "Hello, "

// Generate constructs a localized greeting message for the specified name.
// If the provided name is an empty string, it gracefully defaults to greeting "Go Bank".
func Generate(name string) string {
	if name == "" {
		name = "Go Bank"
	}
	return defaultGreetingPrefix + name + "!"
}
