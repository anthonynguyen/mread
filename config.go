package main

import (
	"os"
)

type Config struct {
	DEBUG bool
	PORT  string
}

var CONFIG Config = Config{
	false,
	"5678",
}

func loadConfig() {
	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "1" {
		CONFIG.DEBUG = false
	}

	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		CONFIG.PORT = portEnv
	}
}
