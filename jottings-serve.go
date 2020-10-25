//go:generate statik -src=./public -include=*.jpg,*.txt,*.html,*.css,*.js

package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	GetTitle() string
	GetText() string
	GetPath() string
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
	SearchIndex     BleveSeachIndex
}

func NewMarkdownIndex(dir string) MarkdownFileIndex {
	files := make(map[string]MarkdownFile)
	searchindex := GetIndex()
	mdfi := MarkdownFileIndex{
		Files:           files,
		Directory:       dir,
		ContentTemplate: contentTemplate(),
		IndexTemplate:   indexTemplate(),
		SearchIndex:     searchindex,
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
		oldfile, exists := mdfi.Files[filename]
		mdfile := ReadMarkdown(path, mdfi.ContentTemplate)
		if exists {
			if oldfile.Checksum != mdfile.Checksum {
				mdfi.SearchIndex.Put(filename, mdfile)
			}
			mdfi.Files[filename] = mdfile
		} else {
			mdfi.SearchIndex.Put(filename, mdfile)
			mdfi.Files[filename] = mdfile
		}
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
	checksum, _ := md5sum(path)
	return MarkdownFile{
		FilePath: path,
		Filename: filename,
		Title:    filename,
		Template: template,
		Checksum: checksum,
	}
}

func md5sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	result := hex.EncodeToString(hash.Sum(nil))
	return result, nil
}

type MarkdownFile struct {
	FilePath string
	Filename string
	Title    string
	Template *template.Template
	Checksum string
}

func (md MarkdownFile) ToHTML(w http.ResponseWriter) {
	fc, err := ioutil.ReadFile(md.FilePath)
	if err != nil {
		log.Printf("Couldn't load file %v", md.FilePath)
	}

	html := string(markdown.ToHTML(fc, nil, nil))

	md.Template.Execute(w, Content{
		Title: md.Title,
		Body:  template.HTML(html),
	})
}

func (md MarkdownFile) GetTitle() string {
	return md.Title
}

func (md MarkdownFile) GetText() string {
	fc, err := ioutil.ReadFile(md.FilePath)
	if err != nil {
		log.Printf(err.Error())
	}

	return string(fc)
}

func (md MarkdownFile) GetPath() string {
	return md.FilePath
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
	input, err := ioutil.ReadFile(mdf.FilePath)
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
	s.Index.ReadFiles()
	file, exists := s.Index.Get(r.URL.Path)
	if exists {
		file.ToHTML(w)
	} else {
		log.Print("This is the index page")
		s.Index.ServeIndex(w)
	}
}
