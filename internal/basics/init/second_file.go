package initializer

// secondFileVar is initialised by this file's own init().
// Go guarantees that within a package, variables are initialised before any init()
// runs, and init() functions run in the order their source files are compiled
// (typically alphabetical). Relying on that ordering is fragile — prefer
// explicit Setup() functions for ordering-sensitive initialisation.
var secondFileVar = "second file default"

func init() {
	secondFileVar = "second file init ran"
}

// GetSecondVar returns the value set by this file's init().
func GetSecondVar() string {
	return secondFileVar
}
