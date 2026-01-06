package main

import (
	"viewer/internal/app"
	"viewer/internal/config"
)

func main() {
	// init config: json
	cfg := config.MustLoad()

	app.Run(cfg)
}
