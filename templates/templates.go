package templates

import (
	_ "embed"
	"html/template"
)

//go:embed index.gohtml
var indexHTML string
var index = template.Must(template.New("index").Parse(indexHTML))

func Index() *template.Template { return index }

type IndexData struct {
	Textarea string
	Error    string
	Result   string
	Link     string
}
