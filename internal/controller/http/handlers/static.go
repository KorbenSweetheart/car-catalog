package handlers

import (
	"log/slog"
	"net/http"
)

func HandleStatic(log *slog.Logger) http.Handler {
	const op = "handlers.static.handle.Static"

	log = log.With(
		slog.String("op", op),
	)

	fs := http.FileServer(http.Dir("static"))

	// 2. Strip the prefix so the file server sees "style.css"
	//    instead of "/static/style.css"
	return http.StripPrefix("/static/", fs)
}
