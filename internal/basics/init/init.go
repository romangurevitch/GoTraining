package initializer

var localVar = "what would be the value?"

func init() {
	localVar = "is that what you expect?"
}

func GetVar() string {
	return localVar
}
