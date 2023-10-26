package monibot

import (
	"os"
	"strings"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestVersion(t *testing.T) {
	ass := assert.New(t)
	// Version must start with "v"
	ass.True(len(Version) >= 6)
	ass.Eq("v", Version[0:1])
	// parse CHANGELOG.md
	data, err := os.ReadFile("CHANGELOG.md")
	ass.Nil(err)
	_, rest, found := strings.Cut(string(data), "## v")
	ass.Eq(true, found)
	rest, _, found = strings.Cut(rest, " ")
	ass.True(found)
	ass.Eq(Version, "v"+rest)
}
