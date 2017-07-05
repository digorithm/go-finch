package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewMember(db *sqlx.DB) *Member {
	member := &Member{}
	member.db = db
	member.table = "member_of"
	member.hasID = false

	return member
}

type Member struct {
	Base
}

//adds to the member_of table the given userID, houseID and
// the own_type corresponding to "owner"
func (m *Member) addOwner(tx *sqlx.Tx, houseID, userID int64) {

	m.addUserHelper(tx, houseID, userID, 1)
}

//adds to the member_of table the given userID, houseID and
// the own_type corresponding to "resident"
func (m *Member) addResident(tx *sqlx.Tx, houseID, userID int64) {

	m.addUserHelper(tx, houseID, userID, 2)
}

// adds to or updates the member_of table with the given userID,
// houseID and the own_type corresponding to "blocked"
func (m *Member) blockUser(tx *sqlx.Tx, houseID, userID int64) {

	m.addUserHelper(tx, houseID, userID, 3)
}

// addUserHelper is the generic add to member_of table
// and each add/block method calls it with
// specific own_type ID
func (m *Member) addUserHelper(tx *sqlx.Tx, houseID, userID, ownID int64) {

	data := make(map[string]interface{})
	data["user_id"] = userID
	data["house_id"] = houseID
	data["own_type"] = ownID

	_, err := m.InsertIntoMultiKeyTable(tx, data)

	if err != nil {

		fmt.Printf("Got error in adding: %v", err)
	}
}
