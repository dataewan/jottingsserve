package main

import (
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

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

type LinksFromFile struct {
	Title string         `json:"title"`
	Links []MarkdownLink `json:"links"`
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

func (mdfi MarkdownFileIndex) GetAllLinks() []LinksFromFile {
	var allLinks []LinksFromFile
	for _, f := range mdfi.Files {
		lff := LinksFromFile{
			Title: f.Title,
			Links: f.GetLinks(),
		}
		allLinks = append(allLinks, lff)
	}
	return allLinks
}

func filterMarkdownLinks(lff []LinksFromFile) []LinksFromFile {
	mdlinks := []LinksFromFile{}
	for _, linksinfile := range lff {
		links := []MarkdownLink{}
		for _, link := range linksinfile.Links {
			if link.LinkType == "Markdown" {
				links = append(links, link)
			}
		}
		mdlinks = append(mdlinks, LinksFromFile{
			Title: linksinfile.Title,
			Links: links,
		})
	}
	return mdlinks
}

func filterBrokenLinks(lff []LinksFromFile, mdfi MarkdownFileIndex) []LinksFromFile {
	mdlinks := []LinksFromFile{}
	for _, linksinfile := range lff {
		links := []MarkdownLink{}
		for _, link := range linksinfile.Links {
			if !mdfi.Exists(link.Destination) {
				links = append(links, link)
			}
		}
		if len(links) > 0 {
			mdlinks = append(mdlinks, LinksFromFile{
				Title: linksinfile.Title,
				Links: links,
			})
		}
	}
	return mdlinks
}

func (mdfi MarkdownFileIndex) CheckBrokenLinks() []LinksFromFile {
	allLinks := mdfi.GetAllLinks()
	mdLinks := filterMarkdownLinks(allLinks)
	brokenLinks := filterBrokenLinks(mdLinks, mdfi)
	return brokenLinks
}
