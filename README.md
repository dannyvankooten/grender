# Grender [![GoDoc](http://godoc.org/github.com/dannyvankooten/grender?status.svg)](http://godoc.org/github.com/dannyvankooten/grender)  [![Build Status](https://travis-ci.org/dannyvankooten/grender.svg)](https://travis-ci.org/dannyvankooten/grender)

Grender is a package that provides functionality for easily rendering HTML templates and JSON or XML data to a HTTP response. It is based on [github.com/unrolled/render](https://github.com/unrolled/render) with some subtle modifications when it comes to rendering HTML templates.

- Templates inheritance: `{{/* extends "master.tmpl" */}}`
- Glob configuration: `templates/*.tmpl`
- Normal templates as partials: `{{ template "footer" .}}`

## Usage
Grender can be used with pretty much any web framework providing you can access the `http.ResponseWriter` from your handler. The rendering functions simply wraps Go's existing functionality for marshaling and rendering data.

- HTML: Uses the [html/template](https://golang.org/pkg/html/template/) package to render HTML templates.
- JSON: Uses the [encoding/json](https://golang.org/pkg/encoding/json/) package to marshal data into a JSON-encoded response.
- XML: Uses the [encoding/xml](https://golang.org/pkg/encoding/xml/) package to marshal data into an XML-encoded response.
- Text: Passes the incoming string straight through to the `http.ResponseWriter`.

```go
// main.go
package main

import (
    "net/http"
    "github.com/dannyvankooten/grender"  
)

func main() {
    r := grender.New(grender.Options{
        Charset: "ISO-8859-1",
        TemplatesGlob: "examples/*.tmpl",
    })
    mux := http.NewServeMux()

    // This will set the Content-Type header to "application/json; charset=ISO-8859-1".
    mux.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
        r.JSON(w, http.StatusOK, map[string]string{"hello": "world"})
    })

    // This will set the Content-Type header to "text/html; charset=ISO-8859-1".
    mux.HandleFunc("/html", func(w http.ResponseWriter, req *http.Request) {
        r.HTML(w, http.StatusOK, "hello.tmpl", "world")
    })

    http.ListenAndServe("127.0.0.1:3000", mux)
}
```

### Options

Grender comes with a variety of configuration options. The defaults are listed below.

```go
r := grender.New(grender.Options{
    Debug: false,       // If true, templates will be recompiled before each render call
    TemplatesGlob: "",  // Glob to your template files
    PartialsGlob: "",   // Glob to your patials or global templates
    Funcs: nil,         // Your template FuncMap
    Charset: "UTF-8",   // Charset to use for Content-Type header values
})
```

### Extending another template

First, define your parent template like this.

file: _master.tmpl_
```go
<html>
  {{template "content" .}}
</html>
```

Then, in a separate template file use a template comment on the first line to indicate that you want to extend the other template file.

file: _child.tmpl_
```go
{{/* extends "master.tmpl" */}}

{{define "content"}}Hello world!{{end}}
```

### More examples

The [grender_test.go](grender_test.go) file contains additional usage examples.

### License

See [LICENSE](LICENSE) file.
