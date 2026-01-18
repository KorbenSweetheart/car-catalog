package webapi

import (
	"log/slog"
	"net/url"
)

func (w *WebRepository) imageURL(carImage string) string {
	const op = "repository.webapi.imageURL"

	log := w.log.With(
		slog.String("op", op),
	)

	if carImage == "" {
		log.Warn("empty car image field")
		return ""
	}

	URL, err := url.JoinPath(w.mediaHost, carImage)
	if err != nil {
		// Fallback: return original if join fails,
		log.Warn("failed to join path",
			slog.String("mediaHost", w.mediaHost),
			slog.String("carImage", carImage),
			slog.Any("error", err),
		)
		return carImage
	}

	return URL
}
