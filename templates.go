package main

import (
	"html/template"
)

func header() string {
	return `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width">
		<style>
				* {
					box-sizing: border-box;
					text-rendering: optimizeLegibility;
					-webkit-font-smoothing: antialiased;
					-moz-osx-font-smoothing: grayscale;
					font-kerning: auto;
					font-family: sans-serif
				}

				body {
					-webkit-text-size-adjust: 100%;
					margin: auto;
					width: 75%;
				}

				img {
					width: 60%;
					margin: auto;
				}
		</style>
    </head>
    <body>
	`
}

func footer() string {
	return `
    </body>
</html>
	`
}

func indexTemplate() *template.Template {
	definition := `
	<h1>Hiya</h1>
	{{range .Pages}}
	<ul>
		<a href={{.Path}}>{{.Filename}}</a>
	</ul>
	{{end}}
	`
	return createTemplate(definition)
}

func contentTemplate() *template.Template {
	definition := `
	<a href="/">home</a>
	<h1>{{.Title}}</h1>
	{{.Body}}
	`

	return createTemplate(definition)
}

func createTemplate(templateDefinition string) *template.Template {
	tmpl, _ := template.New("index").Parse(header() +
		templateDefinition +
		footer())
	return tmpl
}
