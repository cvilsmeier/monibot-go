package api

import (
	"testing"
)

// assertions (maybe use testify or matryer/is when codebase gets larger)

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
