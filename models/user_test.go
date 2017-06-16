package models

import (
	_ "github.com/lib/pq"
	"testing"
)

func newUserForTest(t *testing.T) *User {
	return NewUser(newDbForTest(t))
}

func deleteTestUser(t *testing.T, u *User, id int64) {
	_, err := u.DeleteById(nil, id)
	if err != nil {
		t.Fatal("Something went wrong with test user deletion. Error: %v", err)
	}
}

func TestUserCRUD(t *testing.T) {
	u := newUserForTest(t)

	// Signup
	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
	defer deleteTestUser(t, u, userRow.ID)

	if err != nil {
		t.Errorf("Signing up user should work. Error: %v", err)
	}
	if userRow == nil {
		t.Fatal("Signing up user should work.")
	}
	if userRow.ID <= 0 {
		t.Fatal("Signing up user should work.")
	}
}

func TestGetUserById(t *testing.T) {
	u := newUserForTest(t)

	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
	defer deleteTestUser(t, u, userRow.ID)

	returningUserRow, err := u.GetById(nil, userRow.ID)

	if err != nil {
		t.Errorf("Find user by ID should work")
	}

	if userRow.ID != returningUserRow.ID {
		t.Errorf("IDs did not match!")
	}

}

func TestGetUserByEmail(t *testing.T) {
	u := newUserForTest(t)

	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
	defer deleteTestUser(t, u, userRow.ID)

	returningUserRow, err := u.GetByEmail(nil, userRow.Email)

	if err != nil {
		t.Errorf("Find user by Email should work")
	}

	if userRow.Email != returningUserRow.Email {
		t.Errorf("Emails did not match!")
	}

}

func TestGetUserByUsername(t *testing.T) {
	u := newUserForTest(t)

	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
	defer deleteTestUser(t, u, userRow.ID)

	returningUserRow, err := u.GetByUsername(nil, userRow.Username)

	if err != nil {
		t.Errorf("Find user by Username should work")
	}

	if userRow.Username != returningUserRow.Username {
		t.Errorf("Usernames did not match!")
	}

}
