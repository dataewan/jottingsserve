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
	"github.com/gorilla/mux"
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
		Index:      mdfi,
		FileServer: http.FileServer(http.Dir(Directory)),
	}

	r := mux.NewRouter().StrictSlash(true)

	go browser.OpenURL("http://localhost" + portstring)
	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))
	r.HandleFunc("/", s.HomeHandler)
	r.HandleFunc("/{path}", s.MarkdownFileHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1" + portstring,
	}

	log.Fatal(srv.ListenAndServe())
}

type HTTPServer struct {
	Index      MarkdownFileIndex
	FileServer http.Handler
}

type MarkdownFileIndex struct {
	Files     map[string]MarkdownFile
	Directory string
}

func NewMarkdownIndex(dir string) MarkdownFileIndex {
	files := make(map[string]MarkdownFile)
	mdfi := MarkdownFileIndex{
		Files:     files,
		Directory: dir,
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
		mdfile := ReadMarkdown(path)
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
	template := indexTemplate()
	template.Execute(w, mdfi.Files)
}

func ReadMarkdown(path string) MarkdownFile {
	filename := justFilename(path)
	return MarkdownFile{
		Path:     path,
		Filename: filename,
		Title:    filename,
	}
}

type MarkdownFile struct {
	Path     string
	Filename string
	Title    string
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
	t := contentTemplate()
	t.Execute(w, Content{Title: mdf.Filename, Body: template.HTML(html)})
}

func (s *HTTPServer) MarkdownFileHandler(w http.ResponseWriter, r *http.Request) {
	file, exists := s.Index.Get(r.URL.Path)
	if exists {
		file.ToHTML(w)
	}
}

func (s *HTTPServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	s.Index.ServeIndex(w)
}
