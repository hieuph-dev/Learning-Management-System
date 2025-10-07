package main

import (
	"lms/src/app"
	"lms/src/config"
)

func main() {
	// Initialize configuration
	cfg := config.NewServerConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}
