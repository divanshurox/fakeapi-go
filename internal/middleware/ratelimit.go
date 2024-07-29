package middleware

import (
	"net/http"
	"time"
)

var (
	MaxConnections   = 10
	MaxTimeout       = 1
	TooManyRequests  = "Too many requests"
	RequestCancelled = "Request timed out"
)

var ch = make(chan struct{}, MaxConnections)

func AddRateLimiting(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := time.NewTimer(time.Duration(MaxTimeout) * time.Second)
		ctx := r.Context()
		select {
		case <-timer.C:
			http.Error(w, TooManyRequests, http.StatusTooManyRequests)
			return
		case <-ctx.Done():
			timer.Stop()
			http.Error(w, RequestCancelled, http.StatusGatewayTimeout)
			return
		case ch <- struct{}{}:
			timer.Stop()
			defer func() {
				timer.Stop()
				<-ch
			}()
			handler.ServeHTTP(w, r)
		}
	})
}
