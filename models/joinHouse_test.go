package models

import (
	"testing"

	"fmt"

	_ "github.com/lib/pq"
)

func newJoinForTest(t *testing.T) *Join {
	return NewJoin(newDbForTest(t))
}

func TestGetHouseInvites(t *testing.T) {

	j := newJoinForTest(t)

	result, err := j.GetHouseInvitations(nil, 3)

	if err != nil {
		t.Errorf("Getting house invitations should work. Error: %v", err)
	}

}

func TestAddInvitation(t *testing.T) {

	j := newJoinForTest(t)

	inviteJSON := []byte(`{"house_id": 2, "user_id": 3}`)

	j.addInvitation(nil, inviteJSON)
}
