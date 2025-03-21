package middleware

// MiddlewareConfig holds the configuration for the different middlewares used in the application.
// It includes configurations for CORS (Cross-Origin Resource Sharing) and Logging middleware.
type MiddlewareConfig struct {
	CORSConfig    CORSConfig    // CORS configuration for handling cross-origin requests
	LoggingConfig LoggingConfig // Logging configuration for logging HTTP request details
}
