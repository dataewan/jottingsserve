//go:generate statik -src=./public -include=*.jpg,*.txt,*.html,*.css,*.js

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type Server interface {
	Serve()
}

type File interface {
	ToHTML(http.ResponseWriter)
}

func main() {
	var port string
	var directory string

	app := &cli.App{
		Name:                 "jottingsserve",
		EnableBashCompletion: true,
		Version:              "0.0.1",
		Description:          "Tools for working with markdown linked notes",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Destination: &port,
				Value:       "8080",
			},
			&cli.StringFlag{
				Name:        "directory",
				Aliases:     []string{"d"},
				Destination: &directory,
				Value:       ".",
			},
		},
		Action: func(c *cli.Context) error {
			srv := NewServer(c.String("port"), c.String("directory"))
			log.Fatal(srv.Server.ListenAndServe())
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "checklinks",
				Action: func(c *cli.Context) error {
					TerminalOutputMissingLinks(c.String("directory"))
					return nil
				},
				Usage: "Check for missing links",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
