package grender

import (
	"bufio"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
)

var extendsRegex *regexp.Regexp

func init() {
	var err error
	extendsRegex, err = regexp.Compile(`\{\{\/\* *?extends +?"(.+?)" *?\*\/\}\}`)
	if err != nil {
		panic(err)
	}
}

func (r *Grender) compileTemplatesFromDir() {
	if r.Options.TemplatesGlob == "" {
		return
	}

	// replace existing templates.
	// NOTE: this is unsafe, but Debug should really not be true in production environments.
	templateSet := make(map[string]*template.Template)

	files, err := filepath.Glob(r.Options.TemplatesGlob)
	if err != nil {
		panic(err)
	}

	baseTmpl := template.New("").Funcs(r.Options.Funcs)

	// parse partials (glob)
	if r.Options.PartialsGlob != "" {
		baseTmpl = template.Must(baseTmpl.ParseGlob(r.Options.PartialsGlob))
	}

	for _, templateFile := range files {
		fileName := filepath.Base(templateFile)
		layout := getLayoutForTemplate(templateFile)

		// set template name
		name := fileName
		if layout != "" {
			name = filepath.Base(layout)
		}

		tmpl := template.Must(baseTmpl.Clone())
		tmpl = tmpl.New(name)

		// parse master template
		if layout != "" {
			layoutFile := filepath.Join(filepath.Dir(templateFile), layout)
			tmpl = template.Must(tmpl.ParseFiles(layoutFile))
		}

		// parse child template
		tmpl = template.Must(tmpl.ParseFiles(templateFile))

		templateSet[fileName] = tmpl
	}

	r.Templates.set = templateSet
}

// Lookup returns the compiled template by its filename or nil if there is no such template
func (t *templates) Lookup(name string) *template.Template {
	return t.set[name]
}

// getLayoutForTemplate scans the first line of the template file for the extends keyword
func getLayoutForTemplate(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	b := scanner.Bytes()
	if l := extendsRegex.FindSubmatch(b); l != nil {
		return string(l[1])
	}

	return ""
}
