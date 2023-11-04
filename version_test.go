package monibot

import (
	"os"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestVersion(t *testing.T) {
	// version.go and README.md must have same version
	ass := assert.New(t)
	// Version must start with "v"
	ass.True(len(Version) >= 6)
	ass.Eq("v", Version[0:1])
	// README.md version must be equal
	filename := "README.md"
	data, err := os.ReadFile(filename)
	ass.Nil(err)
	readmeVersion, found := cutout(string(data), "### v", "\n")
	ass.True(found)
	ass.Eq(filename+": "+Version, filename+": v"+readmeVersion)
}
