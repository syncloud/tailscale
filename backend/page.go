package main

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed index.html
var indexHtml string

type Page struct {
	template *template.Template
	url      string
}

func NewPage(url string) *Page {
	return &Page{template: template.Must(template.New("index").Parse(indexHtml)), url: url}
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = p.template.Execute(w, p)
}

func (p *Page) Url() string {
	return p.url
}
