package monibot

import (
	"os"
	"strings"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestExampleSnippets(t *testing.T) {
	ass := assert.New(t)
	// parse example/main.go
	data, err := os.ReadFile("example/main.go")
	ass.Nil(err)
	want, found := cutout(normalizeText(string(data)), "func main() {", "}\n}")
	ass.True(found)
	// parse README.md
	filename := "README.md"
	data, err = os.ReadFile(filename)
	ass.Nil(err)
	have, found := cutout(normalizeText(string(data)), "func main() {", "}\n}")
	ass.True(found)
	ass.Eq(filename+":"+want, filename+":"+have)
	// parse doc.go
	filename = "doc.go"
	data, err = os.ReadFile(filename)
	ass.Nil(err)
	have, found = cutout(normalizeText(string(data)), "func main() {", "}\n}")
	ass.True(found)
	ass.Eq(filename+":"+want, filename+":"+have)
}

func normalizeText(s string) string {
	s = replace(s, "\r", "")
	s = replace(s, "\t", " ")
	s = replace(s, "  ", " ")
	s = replace(s, "\n ", "\n")
	return strings.TrimSpace(s)
}

func replace(str, old, new string) string {
	i := 0
	for strings.Contains(str, old) && i < 1000 {
		str = strings.ReplaceAll(str, old, new)
		i++
	}
	return str
}

func cutout(s, pre, post string) (string, bool) {
	i := strings.Index(s, pre)
	if i < 0 {
		return "", false
	}
	s = s[i+len(pre):]
	i = strings.Index(s, post)
	if i < 0 {
		return "", false
	}
	return s[:i], true
}
