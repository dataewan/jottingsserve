package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

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
	SectionName string `json:"sectionname"`
	SectionHTML string `json:"sectionhtml"`
	SectionRaw  string `json:"sectionraw"`
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

func ReadMarkdown(path string) MarkdownFile {
	filename := justFilename(path)
	return MarkdownFile{
		Path:     path,
		Filename: filename,
		Title:    filename,
	}
}

func singleNodeAsHTML(node ast.Node) string {
	renderer := newNodeRenderer()
	return string(markdown.Render(node, renderer))
}

func singleNodeRawContents(node ast.Node) string {
	fmt.Printf("%+v\n", node)
	str := ""
	return str
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

func splitSections(contents string) []string {
	var sections []string
	section := ""
	re := regexp.MustCompile("^#+")

	sc := bufio.NewScanner(strings.NewReader(contents))

	for sc.Scan() {
		line := sc.Text()
		if re.MatchString(line) {
			sections = append(sections, section)
			section = ""
		}
		section = section + line + "\n"
	}

	sections = append(sections, section)

	return sections
}

func parseSection(contents string) Section {
	var headingName string
	parser := newParser()
	tree := markdown.Parse([]byte(contents), parser)
	html := singleNodeAsHTML(tree)

	for _, node := range tree.GetChildren() {
		if heading, ok := node.(*ast.Heading); ok {
			headingName = (nodeChildText(heading))
		}
	}

	return Section{
		SectionName: headingName,
		SectionHTML: html,
		SectionRaw:  contents,
	}
}

func ParseFileContents(title string, contents []byte) FileContents {
	fc := FileContents{
		Title: title,
	}

	sections := splitSections(string(contents))

	for _, section := range sections {
		fc.Sections = append(fc.Sections, parseSection(section))
	}

	return fc
}

func (mdf MarkdownFile) GetFileContents() FileContents {
	contents := mdf.readFile()
	return ParseFileContents(mdf.Title, contents)
}

type linkVisitor struct {
	nodeQueue  []ast.Node
	withinLink bool
	Links      []MarkdownLink
}

type MarkdownLink struct {
	Destination string `json:"destination"`
	Text        string `json:"text"`
	LinkType    string `json:"linktype"`
}

func (mdl *MarkdownLink) getLinkType() {
	// if mdl.destination ends in markdown and doesn't have a http in it, then markdown
	// if it does have http in it then it is external
	// and also if it has wikipedia in it then it is a wikipedia link
	if strings.HasSuffix(mdl.Destination, ".md") && !strings.Contains(mdl.Destination, "http") {
		mdl.LinkType = "Markdown"
	}

	if strings.Contains(mdl.Destination, "http") {
		if strings.Contains(mdl.Destination, "wikipedia") {
			mdl.LinkType = "Wikipedia"
		} else {
			mdl.LinkType = "External"
		}
	}
}

func (lv *linkVisitor) nodeQueueToLink() MarkdownLink {
	mdl := MarkdownLink{}

	for _, n := range lv.nodeQueue {
		if link, ok := n.(*ast.Link); ok {
			mdl.Destination = string(link.Destination)
		}
		if text, ok := n.(*ast.Text); ok {
			mdl.Text = string(text.Literal)
		}
	}
	lv.nodeQueue = lv.nodeQueue[:0]

	mdl.getLinkType()
	return mdl
}

func (lv *linkVisitor) Visit(n ast.Node, entering bool) ast.WalkStatus {
	if _, ok := n.(*ast.Link); ok {
		if entering {
			lv.nodeQueue = append(lv.nodeQueue, n)
			lv.withinLink = true
		}
		if !entering {
			lv.Links = append(lv.Links, lv.nodeQueueToLink())
			lv.withinLink = false
		}
	}

	if lv.withinLink {
		lv.nodeQueue = append(lv.nodeQueue, n)
	}
	return ast.GoToNext
}

func (mdf MarkdownFile) GetLinks() []MarkdownLink {
	contents := mdf.readFile()
	parser := newParser()
	tree := markdown.Parse([]byte(contents), parser)
	lv := linkVisitor{}

	ast.Walk(tree, &lv)
	return lv.Links

}
