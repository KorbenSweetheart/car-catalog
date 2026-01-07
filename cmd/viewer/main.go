package main

import (
	"os"
	"viewer/internal/app"
	"viewer/internal/config"
)

func main() {
	// init config: json
	cfg := config.MustLoad()

	if err := app.Run(cfg); err != nil {
		os.Exit(1)
	}
}
