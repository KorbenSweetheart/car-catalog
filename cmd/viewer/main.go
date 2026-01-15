package main

import (
	"os"
	"viewer/internal/app"
	"viewer/internal/config"
)

func main() {
	// init config: json
	cfg := config.MustLoad()

	app := app.New(cfg)
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
