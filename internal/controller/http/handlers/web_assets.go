package handlers

import (
	"log/slog"
	"net/http"
	"os"
)

func HandleWebAssets(log *slog.Logger) http.Handler {
	const op = "handlers.web-assets.handle.web.assets"

	log = log.With(
		slog.String("op", op),
	)

	info, err := os.Stat("web") // TODO: move checks to config, and remove panic from here.
	if os.IsNotExist(err) {
		log.Error("static directory does not exist", slog.Any("error", err))
		panic(err)
	}
	if !info.IsDir() {
		log.Error("path is not a directory", slog.Any("error", err))
	}
	if err != nil {
		log.Error("error checking static directory", slog.Any("error", err))
	}

	fs := http.FileServer(http.Dir("web"))

	// 2. Strip the prefix so the file server sees "style.css"
	//    instead of "/web/style.css"
	return http.StripPrefix("/web/", fs)
}
