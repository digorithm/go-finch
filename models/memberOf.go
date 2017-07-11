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
func (m *Member) AddOwner(tx *sqlx.Tx, houseID, userID int64) {

	m.AddUserHelper(tx, houseID, userID, 1)
}

//adds to the member_of table the given userID, houseID and
// the own_type corresponding to "resident"
func (m *Member) AddResident(tx *sqlx.Tx, houseID, userID int64) {

	m.AddUserHelper(tx, houseID, userID, 2)
}

// adds to or updates the member_of table with the given userID,
// houseID and the own_type corresponding to "blocked"
func (m *Member) BlockUser(tx *sqlx.Tx, houseID, userID int64) {

	m.AddUserHelper(tx, houseID, userID, 3)
}

// addUserHelper is the generic add to member_of table
// and each add/block method calls it with
// specific own_type ID
func (m *Member) AddUserHelper(tx *sqlx.Tx, houseID, userID, ownID int64) {

	data := make(map[string]interface{})
	data["user_id"] = userID
	data["house_id"] = houseID
	data["own_type"] = ownID

	_, err := m.InsertIntoMultiKeyTable(tx, data)

	if err != nil {

		fmt.Printf("Got error in adding: %v", err)
	}
}

//deleteResident deletes the resident from the house

func (m *Member) DeleteMember(tx *sqlx.Tx, houseID, userID int64) {

	query := fmt.Sprintf("SELECT HOUSE_ID, USER_ID, OWN_TYPE FROM MEMBER_OF WHERE HOUSE_ID = %v AND USER_ID = $1", houseID)
	data, err := m.GetCompoundModel(tx, query, userID)

	if err != nil {

		fmt.Printf("Got error in DeleteMember: %v", err)
	}

	member := createMemberOfRow(data)

	if member.OwnType == 1 {

		m.DeleteOwner(tx, member.HouseID)

	} else if member.OwnType == 2 {

		m.DeleteResident(tx, member.HouseID, member.UserID)

	} else {

		//Deal with other cases in the future
	}

}

func (m *Member) DeleteResident(tx *sqlx.Tx, houseID, userID int64) {

	where := fmt.Sprintf("house_id = %v and user_id = %v and own_type = 2", houseID, userID)
	_, err := m.DeleteFromTable(tx, where)

	if err != nil {

		fmt.Printf("Got error in deleteresident: %v", err)
	}
}

func (m *Member) DeleteOwner(tx *sqlx.Tx, houseID int64) {

	where := fmt.Sprintf("house_id = %v and own_type = 1", houseID)
	_, err := m.DeleteFromTable(tx, where)

	if err != nil {

		fmt.Printf("Got error in deleteowner: %v", err)
	}
}
