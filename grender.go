package grender

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const (
	// ContentHTML HTTP header value for HTML data
	ContentHTML = "text/html"

	// ContentJSON HTTP header value for JSON data
	ContentJSON = "application/json"

	// ContentType HTTP header name for defining the content type
	ContentType = "Content-Type"

	// DefaultCharset for when no specific Charset options was given
	DefaultCharset = "UTF-8"
)

var extendsRegex *regexp.Regexp

// Grender provides functions for easily writing HTML templates & JSON out to a HTTP Response.
type Grender struct {
	options   Options
	templates map[string]*template.Template
}

// Options holds the configuration options for a Renderer
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

func init() {
	var err error
	extendsRegex, err = regexp.Compile(`{{\/\*\s+extends\s+"(.*)"\s+\*\/}}`)
	if err != nil {
		panic(err)
	}
}

// New creates a new Renderer with the given options
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
		options: opts,
	}

	r.compileTemplatesFromDir()
	return r
}

// HTML executes the template and writes to the responsewriter
func (r *Grender) HTML(w http.ResponseWriter, statusCode int, templateName string, data interface{}) error {
	// re-compile on every render call when Debug is true
	if r.options.Debug {
		r.compileTemplatesFromDir()
	}

	tmpl, ok := r.templates[templateName]
	if !ok {
		return fmt.Errorf("unrecognised template %s", templateName)
	}

	// send response headers + body
	w.Header().Set("Content-Type", ContentHTML+"; charset="+r.options.Charset)
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
	w.Header().Set("Content-Type", ContentJSON+"; charset="+r.options.Charset)

	// do nothing if nil data
	if data == nil {
		return nil
	}

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}

func (r *Grender) compileTemplatesFromDir() {
	if r.options.TemplatesGlob == "" {
		return
	}

	// replace existing templates
	r.templates = make(map[string]*template.Template)

	files, err := filepath.Glob(r.options.TemplatesGlob)
	if err != nil {
		panic(err)
	}

	for _, templateFile := range files {
		fileName := filepath.Base(templateFile)
		layout := r.getLayoutForTemplate(templateFile)

		// set template name
		name := fileName
		if layout != "" {
			name = filepath.Base(layout)
		}

		tmpl := template.New(name).Funcs(r.options.Funcs)

		// parse partials (glob)
		if r.options.PartialsGlob != "" {
			tmpl = template.Must(tmpl.ParseGlob(r.options.PartialsGlob))
		}

		// parse master template
		if layout != "" {
			layoutFile := filepath.Join(filepath.Dir(templateFile), layout)
			tmpl = template.Must(tmpl.ParseFiles(layoutFile))
		}

		// parse child template
		tmpl = template.Must(tmpl.ParseFiles(templateFile))

		r.templates[fileName] = tmpl
	}
}

// getLayoutForTemplate scans the template file for the extends keyword
func (r *Grender) getLayoutForTemplate(filename string) string {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if l := extendsRegex.FindSubmatch(scanner.Bytes()); l != nil {
			return string(l[1])
		}
	}

	return ""
}
