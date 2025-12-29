package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Bytes      int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.StatusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.Bytes += size
	return size, err
}

func NewLoggingMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {

			// Retrieve reqID from Context
			// We expect the RequestID middleware to have put it there.
			var reqID string

			// Try to get reqID from Context first
			if v, ok := r.Context().Value(RequestIDKey).(string); ok {
				reqID = v
			} else {
				reqID = r.Header.Get("X-Request-ID")
				if reqID == "" {
					reqID = "unknown"
				}
			}

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", reqID),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			// Wrap the writer
			rr := &responseRecorder{
				ResponseWriter: w,
				StatusCode:     http.StatusOK, // Default to 200
			}

			t1 := time.Now()

			// Defer the logging until AFTER the next handler returns
			defer func() {
				entry.Info("request completed",
					slog.Int("status", rr.StatusCode),
					slog.Int("bytes", rr.Bytes),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(rr, r)
		}

		return http.HandlerFunc(fn)
	}
}
