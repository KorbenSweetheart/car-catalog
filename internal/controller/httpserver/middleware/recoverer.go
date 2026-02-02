package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func NewRecoveringMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/recoverer"),
		)

		log.Info("recoverer middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {

			// 1. Defer a function to recover from panics
			defer func() {
				if err := recover(); err != nil {
					// 1. Check for ErrAbortHandler (Crucial)
					if err == http.ErrAbortHandler {
						panic(err)
					}

					log.Error("panic recovered",
						slog.Any("error", err),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("stack", string(debug.Stack())), // Crucial for debugging
					)

					// 3. Return a 500 Internal Server Error to the user
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}
			}()

			// 4. Call the next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
