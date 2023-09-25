package monibot

import (
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	// Version must start with "v"
	assertTrue(t, len(Version) >= 6)
	assertEq(t, "v", Version[0:1])
	// parse CHANGELOG.md
	data, err := os.ReadFile("CHANGELOG.md")
	assertNil(t, err)
	_, rest, found := strings.Cut(string(data), "## v")
	assertEq(t, true, found)
	rest, _, found = strings.Cut(rest, " ")
	assertTrue(t, found)
	assertEq(t, Version, "v"+rest)
}
