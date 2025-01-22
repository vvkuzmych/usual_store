package main

import (
	"golang.org/x/time/rate"
	"net/http"
)

var ipLimiter = make(map[string]*rate.Limiter)

func getLimiter(ip string) *rate.Limiter {
	if limiter, exists := ipLimiter[ip]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(1, 5) // 1 request per second, burst of 5
	ipLimiter[ip] = limiter
	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
