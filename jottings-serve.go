//go:generate statik -src=./public -include=*.jpg,*.txt,*.html,*.css,*.js

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/dataewan/jottingsserve/statik"
	"github.com/gomarkdown/markdown"
	"github.com/pkg/browser"
	"github.com/rakyll/statik/fs"
)

type Server interface {
	Serve()
}

type FileIndex interface {
	Get(string) File
}

type File interface {
	ToHTML(http.ResponseWriter)
}

var Directory = "."

func main() {
	port := flag.String("port", "8080", "Port to serve pages on")
	flag.Parse()

	portstring := ":" + *port

	if flag.NArg() == 1 {
		Directory = flag.Arg(0)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	mdfi := NewMarkdownIndex(Directory)
	mdfi.ReadFiles()
	s := HTTPServer{
		Index: mdfi,
	}

	go browser.OpenURL("http://localhost" + portstring)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))
	http.HandleFunc("/", s.Serve)
	log.Fatal(http.ListenAndServe(portstring, nil))
}

type HTTPServer struct {
	Index MarkdownFileIndex
}

type MarkdownFileIndex struct {
	Files           map[string]MarkdownFile
	Directory       string
	ContentTemplate *template.Template
	IndexTemplate   *template.Template
}

func NewMarkdownIndex(dir string) MarkdownFileIndex {
	files := make(map[string]MarkdownFile)
	mdfi := MarkdownFileIndex{
		Files:           files,
		Directory:       dir,
		ContentTemplate: contentTemplate(),
		IndexTemplate:   indexTemplate(),
	}

	return mdfi
}

func (mdfi *MarkdownFileIndex) ReadFiles() {
	matches, err := filepath.Glob(mdfi.Directory + "/*md")
	if err != nil {
		log.Print(err)
	}

	for _, path := range matches {
		filename := justFilename(path)
		mdfile := ReadMarkdown(path, mdfi.ContentTemplate)
		mdfi.Files[filename] = mdfile
	}
}

func (mdfi *MarkdownFileIndex) Get(url string) (File, bool) {
	lookup := justFilename(url)
	value, exists := mdfi.Files[lookup]
	if exists {
		return value, true
	}
	return MarkdownFile{}, false
}

func (mdfi *MarkdownFileIndex) ServeIndex(w http.ResponseWriter) {
	template := mdfi.IndexTemplate
	template.Execute(w, mdfi.Files)
}

func ReadMarkdown(path string, template *template.Template) MarkdownFile {
	filename := justFilename(path)
	return MarkdownFile{
		Path:     path,
		Filename: filename,
		Title:    filename,
		Template: template,
	}
}

type MarkdownFile struct {
	Path     string
	Filename string
	Title    string
	Template *template.Template
}

func (md MarkdownFile) ToHTML(w http.ResponseWriter) {
	fc, err := ioutil.ReadFile(md.Path)
	if err != nil {
		log.Printf("Couldn't load file %v", md.Path)
	}

	html := string(markdown.ToHTML(fc, nil, nil))

	md.Template.Execute(w, Content{
		Title: md.Title,
		Body:  template.HTML(html),
	})
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

func (mdf MarkdownFile) readFile() []byte {
	input, err := ioutil.ReadFile(mdf.Path)
	if err != nil {
		return nil
	}
	return input
}

func markdownToHTML(input []byte) []byte {
	html := markdown.ToHTML(input, nil, nil)
	return html
}

func (mdf MarkdownFile) writePage(w http.ResponseWriter) {
	input := mdf.readFile()
	html := string(markdownToHTML(input))
	mdf.Template.Execute(w, Content{Title: mdf.Filename, Body: template.HTML(html)})
}

func (s *HTTPServer) Serve(w http.ResponseWriter, r *http.Request) {
	file, exists := s.Index.Get(r.URL.Path)
	if exists {
		file.ToHTML(w)
	} else {
		s.Index.ServeIndex(w)
	}
}
