package main

import (
	_ "github.com/dataewan/jottingsserve/statik"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
)

type HTTPServer struct {
	Index      MarkdownFileIndex
	FileServer http.Handler
	Router     *mux.Router
	Server     *http.Server
}

func NewServer(port string, directory string) *HTTPServer {
	mdfi := NewMarkdownIndex(directory)
	mdfi.ReadFiles()
	router := mux.NewRouter()
	s := &HTTPServer{
		Index:      mdfi,
		FileServer: http.FileServer(http.Dir(Directory)),
		Router:     router,
		Server: &http.Server{
			Handler: router,
			Addr:    ":" + port,
		},
	}

	s.addRoutes()

	return s
}

func getStatik() http.FileSystem {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	return statikFS
}

func (s *HTTPServer) addRoutes() {
	statikFS := getStatik()
	s.Router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))
	s.Router.HandleFunc("/", s.HomeHandler)
	s.Router.HandleFunc("/{path}", s.MarkdownFileHandler)
	s.Router.HandleFunc("/api/files", s.ApiFileList)
	s.Router.HandleFunc("/api/files/{title}", s.ApiFile)

}
