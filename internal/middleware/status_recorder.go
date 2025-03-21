package middleware

import (
	"bytes"
	"net/http"
)

// StatusRecorder is a custom wrapper around the http.ResponseWriter.
// It captures the status code and response body, allowing us to log these details.
type StatusRecorder struct {
	http.ResponseWriter // Embeds the original ResponseWriter to delegate methods
	statusCode   int     // Stores the HTTP status code of the response
	err          error   // Stores any error encountered during response writing
	body         *bytes.Buffer // Captures the response body for logging
	headerWritten bool    // Flag to ensure the header is only written once
}

// WriteHeader captures the status code and prevents the header from being written more than once.
// It also writes the header to the underlying ResponseWriter.
func (rec *StatusRecorder) WriteHeader(code int) {
	if rec.headerWritten {
		return
	}
	rec.statusCode = code
	rec.headerWritten = true
	rec.ResponseWriter.WriteHeader(code)
}

// Write captures the response body and writes it to both the ResponseWriter and the body buffer.
// It also stores any error encountered while writing the response.
func (rec *StatusRecorder) Write(b []byte) (int, error) {
	n, err := rec.ResponseWriter.Write(b)
	if err != nil {
		rec.err = err
	}
	rec.body.Write(b)
	return n, err
}

// Flush ensures that if the ResponseWriter implements the http.Flusher interface, it will flush the response.
func (rec *StatusRecorder) Flush() {
	if flusher, ok := rec.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
