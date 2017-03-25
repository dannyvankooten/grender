package grender

import (
	"bufio"
	"html/template"
	"os"
	"path/filepath"
)

func (r *Grender) compileTemplatesFromDir() {
	if r.options.TemplatesGlob == "" {
		return
	}

	// replace existing templates.
	// NOTE: this is unsafe, but Debug should really not be true in production environments.
	r.templates = make(map[string]*template.Template)

	files, err := filepath.Glob(r.options.TemplatesGlob)
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
