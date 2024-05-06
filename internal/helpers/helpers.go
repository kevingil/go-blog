package helpers

import (
	"regexp"
	"text/template"
	"time"

	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Template inline helper functions
func until(n int) []struct{} {
	return make([]struct{}, n)
}

func date(t *time.Time) string {
	return t.Local().Format("January 2, 2006")

}

func shortDate(t *time.Time) string {
	return t.Local().Format("01/02/06")
}

func v() string {
	currentDate := time.Now()
	formattedDate := currentDate.Format("020122")
	return formattedDate
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

var Functions = template.FuncMap{
	"date":      date,
	"shortDate": shortDate,
	"truncate":  truncate,
	"mdToHTML":  mdToHTML,
	"draft":     draft,
	"until":     until,
	"v":         v,
	"sub": func(a, b int) int {
		return a - b
	},
	"add": func(a, b int) int {
		return a + b
	},
}
