package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

	var f MarkdownFile

	for _, file := range s.Index.Files {
		if file.Title == title {
			f = file
		}
	}

	if f == (MarkdownFile{}) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	writeJsonResponse(f, w, r)

}
