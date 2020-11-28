package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (s *HTTPServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	template := indexTemplate()
	template.Execute(w, s.Index.Files)
}

func (s *HTTPServer) MarkdownFileHandler(w http.ResponseWriter, r *http.Request) {
	file, exists := s.Index.Get(r.URL.Path)
	if exists {
		file.ToHTML(w)
	}
}

func writeJsonResponse(rawdata interface{}, w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(rawdata)
	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *HTTPServer) ApiFileList(w http.ResponseWriter, r *http.Request) {
	writeJsonResponse(s.Index.Files, w, r)
}

func (s *HTTPServer) ApiFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	file, exists := s.Index.Get(title)

	if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	writeJsonResponse(file, w, r)

}
