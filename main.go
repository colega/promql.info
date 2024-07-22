package main

import (
	"embed"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/colega/envconfig"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"

	"github.com/colega/promql.info/templates"
)

type config struct {
	Host string `default:"127.0.0.1"`
	Port string `default:"8080"`
}

//go:embed static
var static embed.FS

//go:embed example.test
var example string

func main() {
	var cfg config
	if err := envconfig.Process("APP", &cfg); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePost(w, r)
		} else {
			handleGet(w, r)
		}
	}))

	err := http.ListenAndServe(cfg.Host+":"+cfg.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var query string
	queryBytes, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.FormValue("b64"))))
	if err == nil && len(queryBytes) > 0 {
		query = string(queryBytes)
	} else {
		query = example
	}
	handleQuery(w, query)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the query value
	query := r.FormValue("query")
	if query == "" {
		data := templates.IndexData{
			Textarea: query,
			Error:    "Input can't be empty",
		}
		if err := templates.Index().Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	handleQuery(w, query)
}

func handleQuery(w http.ResponseWriter, query string) {
	errs := run(query)
	result := ""
	if len(errs) == 0 {
		result = "✅ All tests passed"
	}
	link := fmt.Sprintf("https://promql.info/?b64=%s", base64.StdEncoding.EncodeToString([]byte(query)))

	// Render the template
	data := templates.IndexData{
		Textarea: query,
		Error:    strings.Join(errs, "<br>"),
		Result:   result,
		Link:     link,
	}
	if err := templates.Index().Execute(w, data); err != nil {
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
		msg = "❌ " + m[1]
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
