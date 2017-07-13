package models

import (
	"testing"

	_ "github.com/lib/pq"
)

func newUnitForTest(t *testing.T) *Unit {
	return NewUnit(newDbForTest(t))
}

