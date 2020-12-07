package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *HTTPServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	s.Index.ReadFiles()
	template := indexTemplate()
	template.Execute(w, s.Index.Files)
}

func (s *HTTPServer) MarkdownFileHandler(w http.ResponseWriter, r *http.Request) {
	file, exists := s.Index.Get(r.URL.Path)
	if exists {
		file.ToHTML(w)
	} else {
		s.FileServer.ServeHTTP(w, r)
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
	s.Index.ReadFiles()
	writeJsonResponse(s.Index.Files, w, r)
}

func (s *HTTPServer) ApiGetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	file, exists := s.Index.Get(title)

	if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	writeJsonResponse(file, w, r)

}

func (s *HTTPServer) ApiWriteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}

	s.writefile(title, body)
}

func (s *HTTPServer) ApiGetFileSections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	mdfile, exists := s.Index.Get(title)
	if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	fileContents := mdfile.GetFileContents()
	writeJsonResponse(fileContents, w, r)
}
