package models

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func newMemberForTest(t *testing.T) *Member {
	return NewMember(newDbForTest(t))
}

func TestAddOwner(t *testing.T) {
	m := newMemberForTest(t)
	tearDown := TestSetup(t, m.db)
	defer tearDown(t, m.db)

	m.AddOwner(nil, 3, 5)

	where := "HOUSE_ID = 3 AND USER_ID = 5"
	_, err := m.DeleteFromTable(nil, where)

	if err != nil {
		fmt.Printf("TestAddOwner: Delete instance failed: %v", err)
	}
}
