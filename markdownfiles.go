package main

import (
	"github.com/gomarkdown/markdown"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type MarkdownFile struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Title    string `json:"title"`
}

type Content struct {
	Title string
	Body  template.HTML
}

func (md MarkdownFile) ToHTML(w http.ResponseWriter) {
	fc, err := ioutil.ReadFile(md.Path)
	if err != nil {
		log.Printf("Couldn't load file %v", md.Path)
	}

	html := string(markdown.ToHTML(fc, nil, nil))

	t := contentTemplate()
	t.Execute(w, Content{
		Title: md.Title,
		Body:  template.HTML(html),
	})
}

func (mdf MarkdownFile) readFile() []byte {
	input, err := ioutil.ReadFile(mdf.Path)
	if err != nil {
		return nil
	}
	return input
}
func ReadMarkdown(path string) MarkdownFile {
	filename := justFilename(path)
	return MarkdownFile{
		Path:     path,
		Filename: filename,
		Title:    filename,
	}
}
