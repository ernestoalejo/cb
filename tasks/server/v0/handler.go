// Copyright 2013 The Gorilla Authors. All rights reserved.
// Modified for Client Builder.

/*
Package handlers is a collection of handlers for use with Go's net/http package.
*/
package v0

import (
	"log"
	"net/http"
	"time"

	"github.com/ernestokarim/cb/colors"
)

// loggingHandler is the http.Handler implementation for LoggingHandlerTo and its friends
type loggingHandler struct {
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	logger := responseLogger{w: w}
	h.handler.ServeHTTP(&logger, req)
	writeLog(req, t, logger.status, logger.size)
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP status
// code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

// writeLog writes a log entry for req to w in Apache Common Log Format.
// ts is the timestamp with which the entry should be logged.
// status and size are used to provide the response HTTP status and size.
func writeLog(req *http.Request, ts time.Time, status, size int) {
	var color string
	if status == 200 || status == 304 {
		color = colors.Green
	}
	if status == 500 || status == 403 || status == 404 || status == 0 {
		color = colors.Red
	}
	log.Printf("%s[%d] %s %s (%d)%s\n", color, status, req.Method,
		req.RequestURI, size, colors.Reset)
}

// LoggingHandler return a http.Handler that wraps h and logs requests to out in
// Apache Common Log Format (CLF).
//
// See http://httpd.apache.org/docs/2.2/logs.html#common for a description of this format.
//
// LoggingHandler always sets the ident field of the log to -
func LoggingHandler(h http.Handler) http.Handler {
	return loggingHandler{h}
}
