package tools

import (
	"io"
	"net/http"
)

func BasicHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}

func Engine() http.Handler {
	mux := NewLogMux("./log.log", "httplog")
	mux.Handle("/srv", BasicHello)
	return mux
}
