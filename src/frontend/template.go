package frontend

import (
	"html/template"
	"io"
)

type templateType string

var templateLogin templateType = "login.html"
var templateSignup templateType = "signup.html"
var templateProfileEdit templateType = "profile_edit.html"
var templateProfileView templateType = "profile_view.html"

var templates = []templateType{templateLogin, templateSignup, templateProfileEdit, templateProfileView}

type Template struct {
	Tmpl *template.Template
}

func NewTemplate() *Template {
	var files = make([]string, 0, len(templates))
	for _, t := range templates {
		files = append(files, "templates/"+string(t))
	}
	return &Template{
		Tmpl: template.Must(template.ParseFiles(files...)),
	}
}

func (t *Template) Execute(wr io.Writer, tmpl templateType, data any) error {
	return t.Tmpl.ExecuteTemplate(wr, string(tmpl), data)
}
