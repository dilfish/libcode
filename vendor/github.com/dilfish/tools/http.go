// Copyright 2018 Sean.ZH

package tools

import (
	"log"
	"net/http"
)

// LogResponseWriter implement a interface for http
type LogResponseWriter struct {
	ct     []byte
	status int
	w      http.ResponseWriter
}

// HandlerInfo holds method and handle function
type HandlerInfo struct {
	Method  string
	Handler http.HandlerFunc
}

// LogMux is a mux with log
type LogMux struct {
	mux    *http.ServeMux
	lw     *LogResponseWriter
	logger *log.Logger
	mp     map[string]HandlerInfo
}

// NewLogMux create log mux
func NewLogMux(fn, prefix string) *LogMux {
	lm := &LogMux{}
	lm.mux = http.NewServeMux()
	lm.lw = &LogResponseWriter{}
	lm.logger = InitLog(fn, prefix)
	if lm.logger == nil {
		return nil
	}
	lm.mp = make(map[string]HandlerInfo)
	return lm
}

// Header implement response writer
func (lw *LogResponseWriter) Header() http.Header {
	return lw.w.Header()
}

// Write impl response writer
func (lw *LogResponseWriter) Write(bt []byte) (int, error) {
	lw.ct = bt
	return lw.w.Write(bt)
}

// WriteHeader impl response writer
func (lw *LogResponseWriter) WriteHeader(statusCode int) {
	lw.status = statusCode
	lw.w.WriteHeader(statusCode)
}

// GET is api for logmux get
func (l *LogMux) GET(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var hi HandlerInfo
	hi.Handler = handler
	hi.Method = "GET"
	l.mp[pattern] = hi
	l.mux.HandleFunc(pattern, handler)
}

// POST is api for logmux post
func (l *LogMux) POST(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var hi HandlerInfo
	hi.Handler = handler
	hi.Method = "POST"
	l.mp[pattern] = hi
	l.mux.HandleFunc(pattern, handler)
}

// ServeHTTP is logmux service
func (l *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, p := l.mux.Handler(r)
	hi, ok := l.mp[p]
	if ok == false {
		http.NotFound(w, r)
		return
	}
	if r.Method != hi.Method {
		w.Write([]byte("Not allowed"))
		return
	}
	l.lw.w = w
	// limit to 200 bytes
	r.Body = http.MaxBytesReader(w, r.Body, 200)
	hi.Handler(l.lw, r)
	if len(l.lw.ct) > 100 {
		l.logger.Println("sizeof", len(l.lw.ct))
	} else {
		l.logger.Println(r.URL.Path+r.URL.RawQuery, string(l.lw.ct))
	}
}
