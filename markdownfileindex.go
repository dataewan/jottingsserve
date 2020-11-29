package main

import (
	"log"
	"path/filepath"
)

type FileIndex interface {
	Get(string) File
}
type MarkdownFileIndex struct {
	Files     map[string]MarkdownFile `json:"files"`
	Directory string                  `json:"directory"`
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

func (mdfi *MarkdownFileIndex) Get(url string) (MarkdownFile, bool) {
	lookup := justFilename(url)
	value, exists := mdfi.Files[lookup]
	if exists {
		return value, true
	}
	return MarkdownFile{}, false
}

func justFilename(path string) string {
	basepath := filepath.Base(path)
	ext := filepath.Ext(basepath)
	return basepath[0 : len(basepath)-len(ext)]
}
