package middleware

import "net/http"

// CORSConfig is a placeholder struct for configuring CORS (Cross-Origin Resource Sharing) settings.
// It doesn't currently hold any fields but can be extended for more customization in the future.
type CORSConfig struct{}

// Middleware is the main function for applying the CORS policy to incoming HTTP requests.
// It handles OPTIONS requests and allows cross-origin resource sharing.
func (c *CORSConfig) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// enableCors sets the necessary headers in the HTTP response to enable CORS.
// This allows the server to accept requests from different origins.
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
