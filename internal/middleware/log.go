package middleware

import (
	"FakeAPI/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func AddLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.GetLogger().Info(
			"API Info",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
		)
		handler.ServeHTTP(w, r)
	})
}
