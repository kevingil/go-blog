package helpers

import (
	"errors"
	"regexp"
	"time"

	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Template inline helper functions
func Until(n int) []struct{} {
	return make([]struct{}, n)
}

func Date(t *time.Time) string {
	return t.Local().Format("January 2, 2006")

}

func ShortDate(t *time.Time) string {
	return t.Local().Format("01/02/06")
}

func V() string {
	currentDate := time.Now()
	formattedDate := currentDate.Format("020122")
	return formattedDate
}

func MdToHTML(content string) string {
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

func Truncate(s string) string {

	re := regexp.MustCompile("<[^>]*>")
	plainText := re.ReplaceAllString(s, "")

	result := plainText
	if len(plainText) > 126 {
		result = plainText[:160] + ".."
	}

	return result
}

func Draft(i int) bool {
	return i == 1
}

type Stack []interface{}

func (s *Stack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *Stack) Pop() interface{} {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *Stack) Len() int {
	return len(*s)
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// ValidateEmail validates email based on regex format.
func ValidateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("invalid format")
	}

	return nil
}
