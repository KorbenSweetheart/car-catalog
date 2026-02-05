package cookies

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const (
	viewedCarsCookieName = "viewed_cars"
	maxHistorySize       = 31
)

func ViewedCarIDs(r *http.Request, log *slog.Logger) []int {
	const op = "httpserver.cookies.ViewedCarIDs"

	log = log.With(
		slog.String("op", op),
	)

	cookie, err := r.Cookie(viewedCarsCookieName)
	if err != nil || cookie.Value == "" {
		log.Debug("couldn't get cookie", slog.Any("error", err))
		// Error usually means http.ErrNoCookie, so we just return empty history
		return []int{}
	}

	rawCarIDs := strings.Split(cookie.Value, ",")

	// SECURITY: Cap the input size immediately to prevent processing massive headers
	if len(rawCarIDs) > maxHistorySize {
		rawCarIDs = rawCarIDs[:maxHistorySize]
	}

	// Filter and convert to integers
	var carIDs []int
	for i := range rawCarIDs {
		if id, err := strconv.Atoi(rawCarIDs[i]); err == nil && id > 0 {
			carIDs = append(carIDs, id)
		}
	}

	return carIDs
}

// TrackViewedCar updates the viewed_cars cookie by prepending the current car ID.
func TrackViewedCar(w http.ResponseWriter, r *http.Request, carID int, log *slog.Logger) {
	const op = "httpserver.cookies.TrackViewedCar"

	log = log.With(
		slog.String("op", op),
	)

	// 1. Get existing car views history
	var history []string
	if cookie, err := r.Cookie(viewedCarsCookieName); err == nil && cookie.Value != "" {
		log.Debug("retrieved cookie successfully")
		history = strings.Split(cookie.Value, ",")

		// SECURITY: Cap input
		if len(history) > maxHistorySize {
			history = history[:maxHistorySize]
		}
	}

	// 2. Prepend current ID
	currentID := strconv.Itoa(carID)
	newHistory := append([]string{currentID}, history...)

	// 3. Trim to max size again
	if len(newHistory) > maxHistorySize {
		newHistory = newHistory[:maxHistorySize]
	}

	// 4. Set Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     viewedCarsCookieName,
		Value:    strings.Join(newHistory, ","),
		Path:     "/",               // Accessible everywhere
		MaxAge:   30 * 24 * 60 * 60, // 30 days
		HttpOnly: true,              // Security: Not accessible via JS
		SameSite: http.SameSiteLaxMode,
	})

	log.Debug("new cookie has been set successfully")
}
