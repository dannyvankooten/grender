package grender

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUnexistingTemplate(t *testing.T) {
	var err error
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
	})

	w := httptest.NewRecorder()
	err = render.HTML(w, http.StatusOK, string(time.Now().UnixNano())+".tmpl", nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestTemplateExtends(t *testing.T) {
	var err error
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
	})

	w := httptest.NewRecorder()
	err = render.HTML(w, http.StatusOK, "child.tmpl", nil)
	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	expected := "Hello world! from the master template.\n"
	if string(body) != expected {
		t.Errorf("invalid body: expected \"%s\", got \"%s\"", expected, body)
	}
}

func TestTemplatePartial(t *testing.T) {
	var err error
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
		PartialsGlob:  "examples/partials/*.tmpl",
	})

	w := httptest.NewRecorder()
	err = render.HTML(w, http.StatusOK, "child-with-partial.tmpl", nil)
	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	expected := "Hello world! How are we today? from the master template.\n"
	if string(body) != expected {
		t.Errorf("invalid body: expected \"%s\", got \"%s\"", expected, body)
	}
}

// func TestTemplateInSubdirectory(t *testing.T) {
// 	var err error
// 	render := New(Options{
// 		TemplatesGlob: "examples/**/*.tmpl",
// 	})
//
// 	w := httptest.NewRecorder()
// 	err = render.HTML(w, http.StatusOK, "subdir/home.tmpl", nil)
// 	if err != nil {
// 		t.Errorf("expected %#v, got %#v", nil, err)
// 	}
//
// 	res := w.Result()
// 	body, _ := ioutil.ReadAll(res.Body)
// 	expected := "Hello, from a subdirectory.\n"
// 	if string(body) != expected {
// 		t.Errorf("invalid body: expected \"%s\", got \"%s\"", expected, body)
// 	}
// }

func BenchmarkGrenderGetLayoutForFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getLayoutForTemplate("examples/child.tmpl")
	}
}

func BenchmarkGrenderCompileTemplatesFromDir(b *testing.B) {
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
	})

	for i := 0; i < b.N; i++ {
		render.compileTemplatesFromDir()
	}
}
