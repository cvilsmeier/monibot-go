package api

import (
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	// parse CHANGELOG.md
	data, err := os.ReadFile("../CHANGELOG.md")
	assertNil(t, err)
	_, rest, found := strings.Cut(string(data), "## v")
	assertEq(t, true, found)
	rest, _, found = strings.Cut(rest, " ")
	assertEq(t, true, found)
	assertEq(t, Version, "v"+rest)
}
