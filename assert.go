package monibot

import (
	"testing"
)

// assert functions (dont want to bring in testify dependency yet)

func assertTrue(t testing.TB, c bool) {
	if !c {
		t.Helper()
		t.Fatalf("want true but have false")
	}
}

func assertNil(t testing.TB, have any) {
	if have != nil {
		t.Helper()
		t.Fatalf("want nil but have %T(%v)", have, have)
	}
}

func assertEq(t testing.TB, want, have any) {
	if want != have {
		t.Helper()
		t.Fatalf("want %T(%v) but have %T(%v)", want, want, have, have)
	}
}
