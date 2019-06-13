// Copyright 2018 Sean.ZH

package tools

import (
	"io"
	"net/http"
)

// BasicHello is test for http
func BasicHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}

// Engine is http handler
func Engine() http.Handler {
	mux := NewLogMux("./log.log", "httplog")
	mux.GET("/srv", BasicHello)
	return mux
}
