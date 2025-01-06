package monibot

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestReadmeVersion(t *testing.T) {
	ass := assert.New(t)
	// README.md
	data, err := os.ReadFile("README.md")
	ass.Nil(err)
	version, ok := cutout(string(data), "### v", "\n")
	ass.True(ok)
	ass.Eq(Version, "v"+version)
}

func TestFixmeMarkers(t *testing.T) {
	ass := assert.New(t)
	filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		ass.Nil(err)
		if entry.IsDir() || strings.HasPrefix(path, ".git/") || strings.HasSuffix(path, "src_test.go") {
			return nil
		}
		code, err := readCode(path)
		ass.Nil(err)
		if strings.Contains(code, "FIXME") {
			t.Fatalf("found FIXME marker in %s", path)
		}
		if strings.Contains(code, "cvvvv") {
			t.Fatalf("found cvvvv marker in %s", path)
		}
		return nil
	})
}

func TestExamples(t *testing.T) {
	ass := assert.New(t)
	const from = "func main() {"
	const to = "\n}"
	// example/main.go is the canonical example
	code, err := readCode("example/main.go")
	ass.Nil(err)
	canonical, ok := cutout(code, from, to)
	ass.True(ok)
	// README.md
	code, err = readCode("README.md")
	ass.Nil(err)
	snip, ok := cutout(code, from, to)
	ass.True(ok)
	ass.Eq(canonical, snip)
	// doc.go
	code, err = readCode("doc.go")
	ass.Nil(err)
	snip, ok = cutout(code, from, to)
	ass.True(ok)
	ass.Eq(canonical, snip)
}

func readCode(name string) (string, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	s := string(data)
	s = replace(s, "\t", " ")
	s = replace(s, "\r", "")
	s = replace(s, "  ", " ")
	s = replace(s, "\n ", "\n")
	return s, nil
}

func replace(s, from, to string) string {
	for strings.Contains(s, from) {
		s = strings.ReplaceAll(s, from, to)
	}
	return s
}

func cutout(s, from, to string) (string, bool) {
	_, after, ok := strings.Cut(s, from)
	if !ok {
		return s, false
	}
	before, _, ok := strings.Cut(after, to)
	if !ok {
		return s, false
	}
	return before, true
}
