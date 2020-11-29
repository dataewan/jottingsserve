package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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

type FileContents struct {
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

type Section struct {
	SectionName     string `json:"sectionname"`
	SectionContents string `json:"sectioncontents"`
}

func newParser() *parser.Parser {
	exts := parser.CommonExtensions
	p := parser.NewWithExtensions(exts)
	return p
}

func newNodeRenderer() *html.Renderer {
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	r := html.NewRenderer(opts)
	return r
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

func (s *HTTPServer) TitleToPath(title string) string {
	path := filepath.Join(s.Index.Directory, title+".md")
	return path
}

func (s *HTTPServer) writefile(title string, body []byte) {
	path := s.TitleToPath(title)
	err := ioutil.WriteFile(path, body, 0644)
	if err != nil {
		log.Print(err)
	}
}

func ReadMarkdown(path string) MarkdownFile {
	filename := justFilename(path)
	return MarkdownFile{
		Path:     path,
		Filename: filename,
		Title:    filename,
	}
}

func singleNodeRender(node ast.Node) string {
	renderer := newNodeRenderer()
	return string(markdown.Render(node, renderer))
}

func nodeChildText(node ast.Node) string {
	val := ""
	for _, child := range node.GetChildren() {
		if textnode, ok := child.(*ast.Text); ok {
			val += string(textnode.Literal)
		}
	}
	return val
}

func (mdf MarkdownFile) GetFileContents() FileContents {
	fc := FileContents{
		Title: mdf.Title,
	}

	sec := Section{
		SectionName: "main",
	}

	contents := mdf.readFile()
	parser := newParser()
	tree := markdown.Parse(contents, parser)

	for _, node := range tree.GetChildren() {
		if heading, ok := node.(*ast.Heading); ok {
			headingName := (nodeChildText(heading))
			fc.Sections = append(fc.Sections, sec)
			sec = Section{SectionName: headingName}
		}
		sec.SectionContents += singleNodeRender(node)
	}

	fc.Sections = append(fc.Sections, sec)
	return fc
}
