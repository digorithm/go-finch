package models

import (
	"encoding/json"
	"fmt"

	"reflect"

	"github.com/jmoiron/sqlx"
)

type Join struct {
	Base
}

func NewJoin(db *sqlx.DB) *Join {
	join := &Join{}
	join.db = db
	join.table = "join_pending"
	join.hasID = true

	return join
}

type Person struct {
	ID   int64
	Name string
}

func CreatePerson(data []interface{}) []Person {

	var row Person
	var person []Person

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)
		person = append(person, row)
	}

	return person

}

/*func (j *Join) GetUserInvitations(tx *sqlx.Tx, userID int64) ([]byte, error) {
SELECT U.ID, U.USERNAME FROM USER_INFO U INNER JOIN MEMBER_OF M ON U.ID = M.USER_ID WHERE M.HOUSE_ID = 3 AND M.OWN_TYPE = 1
SELECT U.ID, U.USERNAME FROM USER_INFO U INNER JOIN JOIN_PENDING P ON P.USER_ID = U.ID WHERE P.TYPE_ID = 1 AND P.HOUSE_ID = 3
}*/

func (j *Join) GetHouseInvitations(tx *sqlx.Tx, houseID int64) ([]byte, error) {

	queryInviter := "SELECT U.ID, U.USERNAME FROM USER_INFO U INNER JOIN MEMBER_OF M ON U.ID = M.USER_ID WHERE M.HOUSE_ID = $1 AND M.OWN_TYPE = 1"
	queryInvitee := "SELECT U.ID, U.USERNAME FROM USER_INFO U INNER JOIN JOIN_PENDING P ON P.USER_ID = U.ID WHERE P.TYPE_ID = 1 AND P.HOUSE_ID = $1"

	inviterI, err := j.GetCompoundModel(tx, queryInviter, houseID)
	if err != nil {
		fmt.Printf("inviterI: %v", err)
	}

	inviteeI, err := j.GetCompoundModel(tx, queryInvitee, houseID)
	if err != nil {
		fmt.Printf("inviteeI: %v", err)
	}

	inviter := CreatePerson(inviterI)
	invitees := CreatePerson(inviteeI)

	finalJSON := buildInviteJSONResponse(inviter, invitees)

	return finalJSON, err

}

func buildInviteJSONResponse(inviter, invitees []Person) []byte {

	finalUsers := make([]map[string]interface{}, 0, 0)

	for _, invitee := range invitees {

		inviteeJ := make(map[string]interface{})
		inviterJ := make(map[string]interface{})
		finalUser := make(map[string]interface{})

		inviteeJ["id"] = invitee.ID
		inviteeJ["name"] = invitee.Name
		inviterJ["id"] = inviter[0].ID
		inviterJ["name"] = inviter[0].Name
		finalUser["invited_user"] = inviteeJ
		finalUser["invited_by"] = inviterJ

		finalUsers = append(finalUsers, finalUser)
	}

	finalUsersJSON, _ := json.MarshalIndent(finalUsers, "", "	")

	return finalUsersJSON
}
