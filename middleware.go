package apex

// func Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		rw := NewResponseWriter(w)
// 		next.ServeHTTP(rw, r)

// 		totalRequests.WithLabelValues(path).Inc()
// 	})
// }
