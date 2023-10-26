package monibot

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestDemoForDoc(t *testing.T) {
	ass := assert.New(t)
	// parse api_test.go
	data, err := os.ReadFile("doc_test.go")
	ass.Nil(err)
	want, found := cutout(normalizeText(data), "// "+"demo-start", "// "+"demo-end")
	ass.True(found)
	want = strings.TrimSpace(strings.ReplaceAll(want, "\t", ""))
	// parse README.md
	data, err = os.ReadFile("README.md")
	ass.Nil(err)
	have, found := cutout(normalizeText(data), "import \"github.com/cvilsmeier/monibot-go\"", "```")
	ass.True(found)
	have = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(have, "\t", ""), "    ", ""))
	ass.Eq(want, have)
	// parse doc.go
	data, err = os.ReadFile("doc.go")
	ass.Nil(err)
	have, found = cutout(normalizeText(data), "import \"github.com/cvilsmeier/monibot-go\"", "Monibot monitors")
	ass.True(found)
	have = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(have, "\t", ""), "    ", ""))
	ass.Eq(want, have)
}

func normalizeText(data []byte) string {
	// remove tabs
	s := strings.ReplaceAll(string(data), "\t", " ")
	// remove space chains
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	// remove indentations
	for strings.Contains(s, "\n ") {
		s = strings.ReplaceAll(s, "\n ", "\n")
	}
	return s
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

// This code is only here to be copied into README.md and doc.go.
// Do not execute.
func DemoForDoc() {
	// ensure it's never executed
	if 2+2 > 1 {
		panic("do not execute")
	}
	// import "github.com/cvilsmeier/monibot-go"
	// demo-start
	// init the api
	apiKey := os.Getenv("MONIBOT_API_KEY")
	api := NewApi(apiKey)
	// reset a watchdog
	err := api.PostWatchdogReset("2f5f6d47183fdf415a7476837351730c")
	if err != nil {
		log.Fatal(err)
	}
	// demo-end
}
