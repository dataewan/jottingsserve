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
		<link rel="stylesheet" type="text/css" href="/public/jottings.css">
		<link rel="stylesheet"
			  href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/10.1.2/styles/foundation.min.css">
		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/10.1.2/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>
		<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
		<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
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
	{{range $key, $file := . }}
	<ul>
		<a href={{$file.Filename}}>{{$file.Filename}}</a>
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
