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

type HandlerInfo struct {
	Method  string
	Handler http.HandlerFunc
}

type LogMux struct {
	mux    *http.ServeMux
	lw     *LogResponseWriter
	logger *log.Logger
	mp     map[string]HandlerInfo
}

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

func (l *LogMux) GET(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var hi HandlerInfo
	hi.Handler = handler
	hi.Method = "GET"
	l.mp[pattern] = hi
	l.mux.HandleFunc(pattern, handler)
}

func (l *LogMux) POST(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var hi HandlerInfo
	hi.Handler = handler
	hi.Method = "POST"
	l.mp[pattern] = hi
	l.mux.HandleFunc(pattern, handler)
}

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
	hi.Handler(l.lw, r)
	l.logger.Println(r.URL.Path+r.URL.RawQuery, string(l.lw.ct))
}
