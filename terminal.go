package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

func terminalTableMissingLinks(ml []LinksFromFile) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{
		"Missing link in file",
		"Link destination",
	})
	for _, filedescription := range ml {
		for _, link := range filedescription.Links {
			t.AppendRow(table.Row{
				filedescription.Title,
				link.Destination,
			})
		}
	}
	t.Render()
}

func TerminalOutputMissingLinks(dir string) {
	mdfi := NewMarkdownIndex(dir)
	mdfi.ReadFiles()
	broken := mdfi.CheckBrokenLinks()
	if len(broken) > 0 {
		terminalTableMissingLinks(broken)
	}
	color.Green("No broken links")
}
