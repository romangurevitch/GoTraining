//go:build !buildTagsDemo_1 && !buildTagsDemo_2

package buildTags

func WelcomeMessage() string {
	return "Hello, this is the Go Training course! No build tag selected."
}
