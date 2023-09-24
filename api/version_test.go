package api

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
	data, err := os.ReadFile("../CHANGELOG.md")
	assertNil(t, err)
	_, rest, found := strings.Cut(string(data), "## v")
	assertEq(t, true, found)
	rest, _, found = strings.Cut(rest, " ")
	assertTrue(t, found)
	assertEq(t, Version, "v"+rest)
	// parse scripts/dist.sh
	data, err = os.ReadFile("../scripts/dist.sh")
	assertNil(t, err)
	_, rest, found = strings.Cut(string(data), "VERSION=\"")
	assertTrue(t, found)
	rest, _, found = strings.Cut(rest, "\"")
	assertTrue(t, found)
	assertEq(t, Version, rest)
}
