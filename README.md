# Grender [![GoDoc](http://godoc.org/github.com/dannyvankooten/grender?status.svg)](http://godoc.org/github.com/dannyvankooten/grender)

Grender is a package that provides functionality for easily rendering HTML templates and JSON or XML data to a HTTP response. It is based on [github.com/unrolled/render](https://github.com/unrolled/render) with some subtle modifications when it comes to rendering HTML templates.

- Templates can extend other templates using a template comment: `{{/* extends "master.tmpl" */}}`
- Configure template files using a glob string: `templates/*.tmpl`
- Support for partials as normal templates: `{{ template "footer" .}}`

_// child.tmpl_
```html
{{/* extends "master.tmpl" */}}

{{define "content"}}Hello world!{{end}}
```

_// master.tmpl_
```html
{{template "content" .}} from the master template.
```

### License

See [LICENSE](LICENSE) file.
