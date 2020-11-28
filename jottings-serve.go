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

func main() {
	port := flag.String("port", "8080", "Port to serve pages on")
	flag.Parse()

	portstring := *port

	var directory = "."

	if flag.NArg() == 1 {
		directory = flag.Arg(0)
	}

	srv := NewServer(portstring, directory)
	go browser.OpenURL("http://localhost:" + portstring)
	log.Fatal(srv.Server.ListenAndServe())
}
