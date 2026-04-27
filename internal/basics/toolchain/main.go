package main

import (
	"fmt"
	"os"
	// TODO: Uncomment to add configuration support
	// "github.com/spf13/viper"
)

const defaultGreetingPrefix = "Hello, "

func Generate(name string) string {
	if name == "" {
		name = "Go Bank"
	}
	return defaultGreetingPrefix + name + "!"
}

// TODO: Uncomment to load configuration using Viper
// func loadConfig() {
// 	viper.SetConfigName("dev")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath("./config")
//
// 	viper.SetDefault("app_name", "hello")
// 	viper.SetDefault("port", 8080)
// 	viper.SetDefault("log_level", "info")
//
// 	if err := viper.ReadInConfig(); err != nil {
// 		fmt.Printf("No config file found: %v\n", err)
// 	}
// }

func main() {
	// TODO: Uncomment to load configuration
	// loadConfig()
	// fmt.Printf("App: %s, Port: %d, LogLevel: %s\n",
	// 	viper.GetString("app_name"),
	// 	viper.GetInt("port"),
	// 	viper.GetString("log_level"),
	// )

	targetName := ""
	if len(os.Args) > 1 {
		targetName = os.Args[1]
	}

	message := Generate(targetName)
	fmt.Println(message)
}
