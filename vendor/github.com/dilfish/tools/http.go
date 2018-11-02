package tools

import (
	"log"
	"net/http"
)

type LogResponseWriter struct {
	ct     []byte
	status int
	w      http.ResponseWriter
}

type LogMux struct {
	mux    *http.ServeMux
	lw     *LogResponseWriter
	logger *log.Logger
}

func NewLogMux(fn, prefix string) *LogMux {
	lm := &LogMux{}
	lm.mux = http.NewServeMux()
	lm.lw = &LogResponseWriter{}
	lm.logger = InitLog(fn, prefix)
	return lm
}

// implement http.ResponseWriter

func (lw *LogResponseWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *LogResponseWriter) Write(bt []byte) (int, error) {
	lw.ct = bt
	return lw.w.Write(bt)
}

func (lw *LogResponseWriter) WriteHeader(statusCode int) {
	lw.status = statusCode
	lw.w.WriteHeader(statusCode)
}

// implement http.ResponseWriter end

func (l *LogMux) Handle(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.mux.HandleFunc(pattern, handler)
}

func (l *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, _ := l.mux.Handler(r)
	l.lw.w = w
	h.ServeHTTP(l.lw, r)
	l.logger.Println(r.URL.Path+r.URL.RawQuery, string(l.lw.ct))
}
