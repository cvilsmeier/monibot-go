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
	// parse scripts/dist.sh
	data, err = os.ReadFile("../scripts/dist.sh")
	assertNil(t, err)
	_, rest, found = strings.Cut(string(data), "VERSION=\"")
	assertEq(t, true, found)
	rest, _, found = strings.Cut(rest, "\"")
	assertEq(t, true, found)
	assertEq(t, Version, rest)
}
