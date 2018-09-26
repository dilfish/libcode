package main

import (
	"flag"
	"fmt"
	"github.com/dilfish/libcode"
	"github.com/dilfish/tools"
	"io"
	"net/http"
	"strings"
)

var flagP = flag.String("d", "", "decode message")
var glc *libcode.LibCode

func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uri is", r.RequestURI)
	if r.RequestURI == "/e/abc" {
		io.WriteString(w, "aaaa")
		return
	}
	uri := r.RequestURI
	etag := "/e/"
	dtag := "/d/"
	if len(uri) > len(etag) && strings.Index(uri, etag) == 0 {
		io.WriteString(w, glc.Encoder(uri[len(etag):]))
		return
	}
	if len(uri) > len(dtag) && strings.Index(uri, dtag) == 0 {
		ret, err := glc.Decoder(uri[len(dtag):])
		if err != nil {
			io.WriteString(w, err.Error())
			return
		} else {
			io.WriteString(w, ret)
			return
		}
	}
	http.NotFound(w, r)
}

func Engine() http.Handler {
	lc, err := libcode.NewLibCode("core_values.txt", "common_han.txt")
	if err != nil {
		panic(err)
	}
	glc = lc
	mux := tools.NewLogMux("./log.log", "icved")
	mux.Handle("/", Handle)
	return mux
}

func main() {
	flag.Parse()
	port := "1024"
	if *flagP == "" {
		fmt.Println("use default port 1024")
	} else {
		port = *flagP
	}
	mux := Engine()
	http.ListenAndServe(":"+port, mux)
}
