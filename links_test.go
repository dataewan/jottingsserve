package main

import (
	"testing"
)

func TestMissingLink(t *testing.T) {
	files := make(map[string]MarkdownFile)
	files["file1"] = MarkdownFile{
		Path:     "file1.md",
		Filename: "file1.md",
		Title:    "file1",
	}
	lff := []LinksFromFile{
		{
			Title: "file2",
			Links: []MarkdownLink{
				{
					Destination: "file3",
					Text:        "hiya",
					LinkType:    "Markdown",
				},
				{
					Destination: "file1",
					Text:        "hiya",
					LinkType:    "Markdown",
				},
			},
		},
	}
	mdfi := MarkdownFileIndex{
		Files:     files,
		Directory: ".",
	}

	mdlinks := filterMarkdownLinks(lff)
	if len(mdlinks[0].Links) != 2 {
		t.Errorf("was expecting 2 markdown links, got %+v", mdlinks)
	}

	missingLinks := filterBrokenLinks(lff, mdfi)

	if len(missingLinks) != 1 {
		t.Errorf("Was expecting 1 missing link, got %+v", missingLinks)
	}
}
