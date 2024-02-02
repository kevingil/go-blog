package views

import (
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Tmpl is a template.
var Tmpl *template.Template

func until(n int) []struct{} {
	return make([]struct{}, n)
}

func date(t *time.Time) string {
	return t.Local().Format("January 2, 2006 15:04:05")
}

func shortDate(t *time.Time) string {
	return t.Local().Format("January 2, 2006")
}

func mdToHTML(content string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	c := []byte(content)
	var buf bytes.Buffer
	if err := md.Convert(c, &buf); err != nil {
		panic(err)
	}

	return buf.String()
}

func truncate(s string) string {

	re := regexp.MustCompile("<[^>]*>")
	plainText := re.ReplaceAllString(s, "")

	result := plainText
	if len(plainText) > 126 {
		result = plainText[:160] + ".."
	}

	return result
}

func draft(i int) bool {
	return i == 1
}

var functions = template.FuncMap{
	"date":      date,
	"shortDate": shortDate,
	"truncate":  truncate,
	"mdToHTML":  mdToHTML,
	"draft":     draft,
	"until":     until,
	"sub": func(a, b int) int {
		return a - b
	},
	"add": func(a, b int) int {
		return a + b
	},
}

func init() {
	// Direcotries to parse
	dirs := []string{
		"./views/*.gohtml",
		"./views/pages/*.gohtml",
		"./views/forms/*.gohtml",
		"./views/components/*.gohtml"}

	//Create a new Tmpl from all directories
	Tmpl = template.New("").Funcs(functions)
	for _, dir := range dirs {
		files, err := filepath.Glob(dir)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			_, err = Tmpl.ParseFiles(file)
			if err != nil {
				panic(err)
			}
		}
	}
}
