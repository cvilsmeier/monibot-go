package monibot

import (
	"os"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestVersion(t *testing.T) {
	ass := assert.New(t)
	// Version must start with "v"
	ass.True(len(Version) >= 6)
	ass.Eq("v", Version[0:1])
	// CHANGELOG.md version must be equal
	filename := "CHANGELOG.md"
	data, err := os.ReadFile(filename)
	ass.Nil(err)
	changelogVersion, found := cutout(string(data), "## v", " ")
	// _, rest, found := strings.Cut(string(data), "## v")
	// ass.True(found)
	// rest, _, found = strings.Cut(rest, " ")
	ass.True(found)
	ass.Eq(filename+": "+Version, filename+": v"+changelogVersion)
}
