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
