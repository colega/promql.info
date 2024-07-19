package main

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
)

type Data struct {
	Textarea string
	Error    string
	Result   string
}

//go:embed index.gohtml
var indexHTML string
var indexTemplate = template.Must(template.New("index").Parse(indexHTML))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePost(w, r)
		} else {
			handleGet(w, r)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// Render the template
	data := Data{
		Textarea: "",
		Error:    "",
		Result:   "",
	}
	if err := indexTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the query value
	query := r.FormValue("query")
	if query == "" {
		data := Data{
			Textarea: query,
			Error:    "Input can't be empty",
		}
		if err := indexTemplate.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	errs := run(query)
	result := ""
	if len(errs) == 0 {
		result = "All tests passed"
	}

	// Render the template
	data := Data{
		Textarea: query,
		Error:    strings.Join(errs, "<br>"),
		Result:   result,
	}
	if err := indexTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var msgRegex = regexp.MustCompile(`[.\n]*Received unexpected error:[\s\n]+(.*)[\n\s]*$`)

var errFailed = errors.New("failed")

type testingT struct {
	errs []string
}

func (*testingT) FailNow() { panic(errFailed) }
func (t *testingT) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if m := msgRegex.FindStringSubmatch(msg); m != nil {
		msg = m[1]
	}
	t.errs = append(t.errs, msg)
}

func run(input string) (result []string) {
	t := &testingT{}
	defer func() {
		if err := recover(); err != nil && err != errFailed {
			panic(err)
		}
		result = t.errs
	}()
	promqltest.RunTest(t, input, newTestEngine())

	return t.errs
}

func newTestEngine() *promql.Engine {
	return promqltest.NewTestEngine(false, 0, promqltest.DefaultMaxSamplesPerQuery)
}
