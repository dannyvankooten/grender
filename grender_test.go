package grender

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Greeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func TestRendererDefaultCharset(t *testing.T) {
	r := New(Options{})

	if r.options.Charset != DefaultCharset {
		t.Errorf("invalid default charset: expected %#v, got %#v", DefaultCharset, r.options.Charset)
	}
}

func TestRendererJSON(t *testing.T) {
	render := New(Options{
		Charset: "ASCII",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.JSON(w, 299, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.Code != 299 {
		t.Errorf("invalid status code: expected %#v, got %#v", 299, res.Code)
	}

	e := ContentJSON + "; charset=" + render.options.Charset
	if v := res.Header().Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	if v := res.Body.String(); v != "{\"one\":\"hello\",\"two\":\"world\"}\n" {
		t.Errorf("invalid response body: expected %#v, got %#v", "{\"one\":\"hello\",\"two\":\"world\"}\n", v)
	}
}

func TestRendererHTML(t *testing.T) {
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "hello.tmpl", "world")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.Code != http.StatusOK {
		t.Errorf("invalid status code: expected %#v, got %#v", http.StatusOK, res.Code)
	}

	e := ContentHTML + "; charset=" + render.options.Charset
	if v := res.Header().Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	if v := res.Body.String(); v != "Hello world!\n" {
		t.Errorf("invalid body: expected %#v, got %#v", "Hello world!\n", v)
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

func BenchmarkGrenderGetLayoutForFile(b *testing.B) {
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
	})

	for i := 0; i < b.N; i++ {
		render.getLayoutForTemplate("examples/child.tmpl")
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
