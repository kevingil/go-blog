package templates

import (
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Tmpl is a template.
var Tmpl *template.Template

func date(t *time.Time) string {
	return t.Local().Format("January 2, 2006 15:04:05")
}

func shortDate(t *time.Time) string {
	return t.Local().Format("January 2, 2006")
}

func mdToHTML(content string) string {
	c := []byte(content)
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(c)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlString := string(markdown.Render(doc, renderer))

	return htmlString
}

func truncate(s string) string {

	re := regexp.MustCompile("<[^>]*>")
	plainText := re.ReplaceAllString(s, "")

	result := plainText
	if len(plainText) > 160 {
		result = plainText[:160] + "..."
	}

	return result
}

var functions = template.FuncMap{
	"date":      date,
	"shortDate": shortDate,
	"truncate":  truncate,
	"mdToHTML":  mdToHTML,
}

func init() {
	Tmpl = template.Must(template.New("./templates/*.html").Funcs(functions).ParseGlob("./templates/*.html"))
}

func init() {
	// Direcotries to parse
	dirs := []string{"./templates/*.html", "./templates/pages/*.html", "./templates/components/*.html"}

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
