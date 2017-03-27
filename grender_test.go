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

func TestGrenderDefaultCharset(t *testing.T) {
	r := New(Options{})

	if r.Options.Charset != DefaultCharset {
		t.Errorf("invalid default charset: expected %#v, got %#v", DefaultCharset, r.Options.Charset)
	}
}

func TestGrenderJSON(t *testing.T) {
	render := New(Options{
		Charset: "ASCII",
	})

	w := httptest.NewRecorder()
	err := render.JSON(w, 299, Greeting{"hello", "world"})
	res := w.Result()

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != 299 {
		t.Errorf("invalid status code: expected %#v, got %#v", 299, res.StatusCode)
	}

	e := ContentJSON + "; charset=" + render.Options.Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if v := string(body); v != "{\"one\":\"hello\",\"two\":\"world\"}\n" {
		t.Errorf("invalid response body: expected %#v, got %#v", "{\"one\":\"hello\",\"two\":\"world\"}\n", v)
	}
}

func TestGrenderHTML(t *testing.T) {
	render := New(Options{
		TemplatesGlob: "examples/*.tmpl",
	})

	w := httptest.NewRecorder()
	err := render.HTML(w, http.StatusOK, "hello.tmpl", "world")
	res := w.Result()
	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("invalid status code: expected %#v, got %#v", http.StatusOK, res.StatusCode)
	}

	e := ContentHTML + "; charset=" + render.Options.Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if v := string(body); v != "Hello world!\n" {
		t.Errorf("invalid body: expected %#v, got %#v", "Hello world!\n", v)
	}
}

func TestGrenderXML(t *testing.T) {
	render := New()

	w := httptest.NewRecorder()
	err := render.XML(w, 299, Greeting{"hello", "world"})
	res := w.Result()

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != 299 {
		t.Errorf("invalid status code: expected %#v, got %#v", 299, res.StatusCode)
	}

	e := ContentXML + "; charset=" + render.Options.Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	expected := "<Greeting><One>hello</One><Two>world</Two></Greeting>"
	if v := string(body); v != expected {
		t.Errorf("invalid response body: expected %#v, got %#v", expected, v)
	}
}

func TestGrenderText(t *testing.T) {
	render := New()

	w := httptest.NewRecorder()
	err := render.Text(w, 200, "Hello world!")
	res := w.Result()

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != 200 {
		t.Errorf("invalid status code: expected %#v, got %#v", 200, res.StatusCode)
	}

	e := ContentText + "; charset=" + render.Options.Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	expected := "Hello world!"
	if v := string(body); v != expected {
		t.Errorf("invalid response body: expected %#v, got %#v", expected, v)
	}
}
