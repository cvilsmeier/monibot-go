package monibot

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestDemoForDoc(t *testing.T) {
	// parse api_test.go
	data, err := os.ReadFile("doc_test.go")
	assertNil(t, err)
	want := normalizeText(data)
	_, want, found := strings.Cut(want, "// "+"demo-start")
	assertTrue(t, found)
	want, _, found = strings.Cut(want, "// "+"demo-end")
	assertTrue(t, found)
	want = strings.ReplaceAll(want, "\t", "")
	want = strings.TrimSpace(want)
	// parse README.md
	data, err = os.ReadFile("README.md")
	assertNil(t, err)
	have := normalizeText(data)
	_, have, found = strings.Cut(have, "import \"github.com/cvilsmeier/monibot-go\"")
	assertTrue(t, found)
	have, _, found = strings.Cut(have, "```")
	assertTrue(t, found)
	have = strings.ReplaceAll(have, "\t", "")
	have = strings.ReplaceAll(have, "    ", "")
	have = strings.TrimSpace(have)
	if want != have {
		t.Logf("README.md: want %q", want)
		t.Logf("README.md: have %q", have)
	}
	assertEq(t, want, have)
	// parse doc.go
	data, err = os.ReadFile("doc.go")
	assertNil(t, err)
	have = normalizeText(data)
	_, have, found = strings.Cut(have, "import \"github.com/cvilsmeier/monibot-go\"")
	assertTrue(t, found)
	have, _, found = strings.Cut(have, "Monibot monitors")
	assertTrue(t, found)
	have = strings.ReplaceAll(have, "\t", "")
	have = strings.ReplaceAll(have, "    ", "")
	have = strings.TrimSpace(have)
	if want != have {
		t.Logf("doc.go: want %q", want)
		t.Logf("doc.go: have %q", have)
	}
	assertEq(t, want, have)
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

// This code is only here to be copied into README.md and doc.go - do not execute.
func DemoForDoc() {
	// ensure it's never executed
	if 2+2 > 1 {
		panic("do not execute")
	}
	// import "github.com/cvilsmeier/monibot-go"
	// demo-start
	// init the api
	userAgent := "my-app/v1.0.0"
	apiKey := os.Getenv("MONIBOT_API_KEY")
	api := NewDefaultApi(userAgent, apiKey)
	// reset a watchdog
	err := api.PostWatchdogReset("000000000000001")
	if err != nil {
		log.Fatal(err)
	}
	// demo-end
}
