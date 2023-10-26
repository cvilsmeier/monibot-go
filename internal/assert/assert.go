// Package assert provides assertion functions for unit tests.
// Its less capable that stretchr/testify but enough for us now.
package assert

import "testing"

type Assert struct {
	t testing.TB
}

func New(t testing.TB) *Assert {
	return &Assert{t}
}

func (a *Assert) True(c bool) {
	a.t.Helper()
	if !c {
		a.t.Fatalf("want true but have false")
	}
}

func (a *Assert) Nil(have any) {
	a.t.Helper()
	if have != nil {
		a.t.Fatalf("want nil but have %T(%v)", have, have)
	}
}

func (a *Assert) Eq(want, have any) {
	a.t.Helper()
	if want != have {
		a.t.Fatalf("want %T(%v) but have %T(%v)", want, want, have, have)
	}
}
