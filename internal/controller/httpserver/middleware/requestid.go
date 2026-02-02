package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type key int

const RequestIDKey key = 0

func NewReqIDMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/requestID"),
		)

		log.Info("requestID middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {

			var sb strings.Builder

			sb.WriteString("req-")
			sb.WriteString(generateRandomID()) // generate custom ID
			reqID := sb.String()

			// Set response header
			w.Header().Set("X-Request-ID", reqID)

			// Set context
			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func generateRandomID() string {
	b := make([]byte, 16) // 16 bytes
	_, err := rand.Read(b)
	if err != nil {
		// In the extremely rare case crypto/rand fails, fallback to a timestamp
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
