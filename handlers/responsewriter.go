package handlers

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		w,
		http.StatusOK,
		0,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	rw.responseSize = len(b)
	return rw.ResponseWriter.Write(b)
}
