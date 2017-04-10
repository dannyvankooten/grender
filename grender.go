package grender

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

const (
	// ContentHTML HTTP header value for HTML data
	ContentHTML = "text/html"

	// ContentJSON HTTP header value for JSON data
	ContentJSON = "application/json"

	// ContentType HTTP header name for defining the content type
	ContentType = "Content-Type"

	// ContentText header value for Text data.
	ContentText = "text/plain"

	// ContentXML header value for XML data.
	ContentXML = "text/xml"

	// DefaultCharset for when no specific Charset Options was given
	DefaultCharset = "UTF-8"
)

type templates struct {
	set map[string]*template.Template
}

// Grender provides functions for easily writing HTML templates & JSON out to a HTTP Response.
type Grender struct {
	Options   Options
	Templates templates
}

// Options holds the configuration Options for a Renderer
type Options struct {
	// With Debug set to true, templates will be recompiled before every render call.
	Debug bool

	// The glob string to your templates
	TemplatesGlob string

	// The Glob string for additional templates
	PartialsGlob string

	// The function map to pass to each HTML template
	Funcs template.FuncMap

	// Charset for responses
	Charset string
}

// New creates a new Renderer with the given Options
func New(optsarg ...Options) *Grender {
	var opts Options

	if len(optsarg) > 0 {
		opts = optsarg[0]
	} else {
		opts = Options{}
	}

	if opts.Charset == "" {
		opts.Charset = "UTF-8"
	}

	r := &Grender{
		Options: opts,
	}

	r.compileTemplatesFromDir()
	return r
}

// HTML executes the template and writes to the responsewriter
func (r *Grender) HTML(w http.ResponseWriter, statusCode int, templateName string, data interface{}) error {
	// re-compile on every render call when Debug is true
	if r.Options.Debug {
		r.compileTemplatesFromDir()
	}

	tmpl := r.Templates.Lookup(templateName)
	if tmpl == nil {
		return fmt.Errorf("unrecognised template %s", templateName)
	}

	// send response headers + body
	w.Header().Set("Content-Type", ContentHTML+"; charset="+r.Options.Charset)
	out := bufPool.Get()
	defer bufPool.Put(out)

	// execute template
	err := tmpl.Execute(out, data)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)
	out.WriteTo(w)
	return nil
}

// JSON renders the data as a JSON HTTP response to the ResponseWriter
func (r *Grender) JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", ContentJSON+"; charset="+r.Options.Charset)

	// do nothing if nil data
	if data == nil {
		return nil
	}

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}

// XML writes the data as a XML HTTP response to the ResponseWriter
func (r *Grender) XML(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", ContentXML+"; charset="+r.Options.Charset)

	// do nothing if nil data
	if data == nil {
		return nil
	}

	w.WriteHeader(statusCode)
	err := xml.NewEncoder(w).Encode(data)
	return err
}

// Text writes the data as a JSON HTTP response to the ResponseWriter
func (r *Grender) Text(w http.ResponseWriter, statusCode int, data string) error {
	w.Header().Set("Content-Type", ContentText+"; charset="+r.Options.Charset)
	w.WriteHeader(statusCode)
	w.Write([]byte(data))
	return nil
}
