package main

import (
	"github.com/gomarkdown/markdown"
	//"github.com/gomarkdown/markdown/html"
	//"github.com/gomarkdown/markdown/parser"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

const INDEXPAGE = "INDEX"
const OTHERPAGE = "OTHER"

type PagePointer struct {
	Path     string
	Filename string
}

type Pages struct {
	Pages []PagePointer
}

type Content struct {
	Title string
	Body  template.HTML
}

func justFilename(path string) string {
	basepath := filepath.Base(path)
	ext := filepath.Ext(basepath)
	return basepath[0 : len(basepath)-len(ext)]
}

func getFiles() Pages {
	matches, err := filepath.Glob("./*md")
	if err != nil {
		log.Print(err.Error())
	}

	var output []PagePointer

	for _, match := range matches {
		filename := justFilename(match)
		output = append(output, PagePointer{Path: match, Filename: filename})
	}

	return Pages{Pages: output}
}

func readFile(page PagePointer) []byte {
	input, err := ioutil.ReadFile(page.Path)
	if err != nil {
		log.Print(err.Error())
	}
	return input
}

func markdownToHTML(input []byte) []byte {
	html := markdown.ToHTML(input, nil, nil)
	return html
}

func checkResponseType(url string) string {
	if url == "/" {
		return INDEXPAGE
	} else {
		return OTHERPAGE
	}
}

func fileIsIndex(file PagePointer) bool {
	if file.Filename == "README" {
		return true
	}
	return false
}

func checkIndexExists() (bool, PagePointer) {
	for _, file := range getFiles().Pages {
		if fileIsIndex(file) {
			return true, file
		}
	}
	return false, PagePointer{}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	exists, page := checkIndexExists()
	if exists {
		writePage(w, page)
	} else {
		markdownfiles := getFiles()
		tmpl := indexTemplate()
		tmpl.Execute(w, markdownfiles)
	}
}

func writePage(w http.ResponseWriter, p PagePointer) {
	input := readFile(p)
	html := markdownToHTML(input)
	tmpl := contentTemplate()
	tmpl.Execute(w, Content{Title: p.Filename, Body: template.HTML(string(html))})
}

func otherpage(w http.ResponseWriter, url string, r *http.Request) {
	markdownfiles := getFiles()

	var markdownfile PagePointer
	found := false

	trimmedURL := justFilename(url)
	for _, file := range markdownfiles.Pages {
		if file.Filename == trimmedURL {
			markdownfile = file
			found = true
		}
	}

	writePage(w, markdownfile)
	if !found {
		http.NotFound(w, r)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	responseType := checkResponseType(r.URL.Path)
	if responseType == INDEXPAGE {
		indexPage(w, r)
	}
	if responseType == OTHERPAGE {
		otherpage(w, r.URL.Path, r)
	}
}

func main() {
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
