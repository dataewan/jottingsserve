//go:generate statik -src=./public -include=*.jpg,*.txt,*.html,*.css,*.js

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pkg/browser"
)

type Server interface {
	Serve()
}

type File interface {
	ToHTML(http.ResponseWriter)
}

var Directory = "."

func main() {
	port := flag.String("port", "8080", "Port to serve pages on")
	flag.Parse()

	portstring := *port

	if flag.NArg() == 1 {
		Directory = flag.Arg(0)
	}

	srv := NewServer(portstring, Directory)
	go browser.OpenURL("http://localhost:" + portstring)
	log.Fatal(srv.Server.ListenAndServe())
}
