package grender

import (
	"bufio"
	"html/template"
	"os"
	"path/filepath"
)

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

	for _, templateFile := range files {
		fileName := filepath.Base(templateFile)
		layout := getLayoutForTemplate(templateFile)

		// set template name
		name := fileName
		if layout != "" {
			name = filepath.Base(layout)
		}

		tmpl := template.New(name).Funcs(r.Options.Funcs)

		// parse partials (glob)
		if r.Options.PartialsGlob != "" {
			tmpl = template.Must(tmpl.ParseGlob(r.Options.PartialsGlob))
		}

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

// getLayoutForTemplate scans the template file for the extends keyword
func getLayoutForTemplate(filename string) string {
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
