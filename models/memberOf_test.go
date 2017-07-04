package models

import (
	"testing"

	_ "github.com/lib/pq"
)

func newMemberForTest(t *testing.T) *Member {
	return NewMember(newDbForTest(t))
}

func TestAddOwner(t *testing.T) {
	m := newMemberForTest(t)

	m.addOwner(nil, 3, 4)
}
