package models

import (
	"encoding/json"
	"fmt"

	"reflect"

	"database/sql"

	"github.com/buger/jsonparser"
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

func (j *Join) GetUserRequests(tx *sqlx.Tx, userID, reqType int64) ([]byte, error) {

	if reqType == 1 {

		return j.GetUserInvitations(tx, userID)
	} else {

		return j.GetUserJoins(tx, userID)
	}
}

func (j *Join) GetUserInvitations(tx *sqlx.Tx, userID int64) ([]byte, error) {

	// queryInvitee gets all the invites
	query := "SELECT U1.ID, U1.USERNAME, U2.ID, U2.USERNAME, H.ID, H.NAME FROM USER_INFO U1 INNER JOIN JOIN_PENDING P ON P.USER_ID = U1.ID INNER JOIN HOUSE H ON P.HOUSE_ID = H.ID INNER JOIN MEMBER_OF M ON M.HOUSE_ID = H.ID INNER JOIN USER_INFO U2 ON M.USER_ID = U2.ID WHERE P.USER_ID = $1 AND M.OWN_TYPE = 1 AND P.TYPE_ID = 1"

	userI, err := j.GetCompoundModel(tx, query, userID)
	if err != nil {
		fmt.Printf("userI: %v", err)
	}

	finalJSON := buildUserInviteJSONResponse(userI)

	return finalJSON, err
}

func (j *Join) GetUserJoins(tx *sqlx.Tx, userID int64) ([]byte, error) {

	// queryInvitee gets all the invites
	query := "SELECT U1.ID, U1.USERNAME, U2.ID, U2.USERNAME, H.ID, H.NAME FROM USER_INFO U1 INNER JOIN JOIN_PENDING P ON P.USER_ID = U1.ID INNER JOIN HOUSE H ON P.HOUSE_ID = H.ID INNER JOIN MEMBER_OF M ON M.HOUSE_ID = H.ID INNER JOIN USER_INFO U2 ON M.USER_ID = U2.ID WHERE P.USER_ID = $1 AND M.OWN_TYPE = 1 AND P.TYPE_ID = 2"

	userI, err := j.GetCompoundModel(tx, query, userID)
	if err != nil {
		fmt.Printf("userI: %v", err)
	}

	finalJSON := buildUserJoinJSONResponse(userI)

	return finalJSON, err
}

func buildUserInviteJSONResponse(userI []interface{}) []byte {

	finalUsers := make([]map[string]interface{}, 0, 0)

	for _, user := range userI {

		v := reflect.ValueOf(user)

		invitee := make(map[string]interface{})
		inviter := make(map[string]interface{})
		house := make(map[string]interface{})
		finalGroup := make(map[string]interface{})

		invitee["id"] = v.Index(0).Interface().(int64)
		invitee["name"] = v.Index(1).Interface().(string)
		inviter["id"] = v.Index(2).Interface().(int64)
		inviter["name"] = v.Index(3).Interface().(string)
		house["house_id"] = v.Index(4).Interface().(int64)
		house["house_name"] = v.Index(5).Interface().(string)
		finalGroup["invited_user"] = invitee
		finalGroup["invited_by"] = inviter
		finalGroup["invited_to"] = house

		finalUsers = append(finalUsers, finalGroup)
	}

	finalUsersJSON, _ := json.MarshalIndent(finalUsers, "", "	")

	return finalUsersJSON
}

func buildUserJoinJSONResponse(userI []interface{}) []byte {

	finalUsers := make([]map[string]interface{}, 0, 0)

	for _, user := range userI {

		v := reflect.ValueOf(user)

		requester := make(map[string]interface{})
		requestingTo := make(map[string]interface{})
		house := make(map[string]interface{})
		finalGroup := make(map[string]interface{})

		requester["id"] = v.Index(0).Interface().(int64)
		requester["name"] = v.Index(1).Interface().(string)
		requestingTo["id"] = v.Index(2).Interface().(int64)
		requestingTo["name"] = v.Index(3).Interface().(string)
		house["house_id"] = v.Index(4).Interface().(int64)
		house["house_name"] = v.Index(5).Interface().(string)
		finalGroup["requesting_user"] = requester
		finalGroup["requested_house"] = house
		finalGroup["house_owner"] = requestingTo

		finalUsers = append(finalUsers, finalGroup)
	}

	finalUsersJSON, _ := json.MarshalIndent(finalUsers, "", "	")

	return finalUsersJSON
}

func (j *Join) GetHouseRequests(tx *sqlx.Tx, houseID, reqType int64) ([]byte, error) {

	if reqType == 1 {

		return j.GetHouseInvitations(tx, houseID)
	} else {

		return j.GetHouseJoins(tx, houseID)
	}
}

func (j *Join) GetHouseJoins(tx *sqlx.Tx, houseID int64) ([]byte, error) {

	queryInviter := "SELECT U.ID, U.USERNAME FROM USER_INFO U INNER JOIN MEMBER_OF M ON U.ID = M.USER_ID WHERE M.HOUSE_ID = $1 AND M.OWN_TYPE = 1"
	queryInvitee := "SELECT U.ID, U.USERNAME FROM USER_INFO U INNER JOIN JOIN_PENDING P ON P.USER_ID = U.ID WHERE P.TYPE_ID = 2 AND P.HOUSE_ID = $1"

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

	finalJSON := buildHouseJoinJSONResponse(inviter, invitees)

	return finalJSON, err
}

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

	finalJSON := buildHouseInviteJSONResponse(inviter, invitees)

	return finalJSON, err

}

func buildHouseInviteJSONResponse(inviter, invitees []Person) []byte {

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

func buildHouseJoinJSONResponse(inviter, invitees []Person) []byte {

	finalUsers := make([]map[string]interface{}, 0, 0)

	for _, invitee := range invitees {

		inviteeJ := make(map[string]interface{})
		inviterJ := make(map[string]interface{})
		finalUser := make(map[string]interface{})

		inviteeJ["id"] = invitee.ID
		inviteeJ["name"] = invitee.Name
		inviterJ["id"] = inviter[0].ID
		inviterJ["name"] = inviter[0].Name
		finalUser["requesting_user"] = inviteeJ
		finalUser["request_to"] = inviterJ

		finalUsers = append(finalUsers, finalUser)
	}

	finalUsersJSON, _ := json.MarshalIndent(finalUsers, "", "	")

	return finalUsersJSON
}

func (j *Join) AddInvitation(tx *sqlx.Tx, inviteJSON []byte) ([]byte, error) {

	inviteEntry := make(map[string]interface{})

	inviteEntry["type_id"] = 1
	inviteEntry["house_id"], _ = jsonparser.GetInt(inviteJSON, "house_id")
	inviteEntry["user_id"], _ = jsonparser.GetInt(inviteJSON, "user_id")

	res, err := j.InsertIntoTable(tx, inviteEntry)
	if err != nil {
		fmt.Printf("addInvite fail: %v", err)
	}

	fmt.Printf("res: %v", res)
	resultJSON := buildAddInviteResponseJSON(res)

	return resultJSON, err

}

func (j *Join) AddJoinRequest(tx *sqlx.Tx, inviteJSON []byte) ([]byte, error) {

	inviteEntry := make(map[string]interface{})

	inviteEntry["type_id"] = 2
	inviteEntry["house_id"], _ = jsonparser.GetInt(inviteJSON, "house_id")
	inviteEntry["user_id"], _ = jsonparser.GetInt(inviteJSON, "user_id")

	res, err := j.InsertIntoTable(tx, inviteEntry)

	resultJSON := buildAddJoinResponseJSON(res)

	return resultJSON, err

}

func (j *Join) FinalizeResponse(tx *sqlx.Tx, responseJSON []byte) ([]byte, error) {

	response := make(map[string]interface{})

	response["accepts"], _ = jsonparser.GetBoolean(responseJSON, "accepts")
	response["invite_id"], _ = jsonparser.GetInt(responseJSON, "invite_id")
	accepts := response["accepts"]
	ID := response["invite_id"]

	var resp []byte
	var err error

	if accepts.(bool) {

		query := "SELECT P.HOUSE_ID, P.USER_ID FROM JOIN_PENDING P WHERE ID = $1"
		res, err := j.GetCompoundModel(tx, query, ID.(int64))

		if err != nil {
			fmt.Printf("getInvitation: %v", err)
		}

		v := reflect.ValueOf(res[0])
		houseID := v.Index(0).Interface().(int64)
		userID := v.Index(1).Interface().(int64)

		m := NewMember(j.db)
		h := NewHouse(j.db)
		m.AddResident(tx, houseID, userID)

		resp, err = h.GetFullHouseInformation(nil, houseID)

	} else {

		resp = []byte(`{"message":"user rejected invitation"}`)

	}

	j.DeleteInvitation(tx, ID.(int64))

	return resp, err
}

func buildAcceptInviteJSON(users []HouseUserOwnRow) []byte {

	finalUser := make(map[string]interface{})
	finalUser["id"] = users[0].HID
	finalUser["household_number"] = users[0].HouseNumber
	members := make([]map[string]interface{}, 0, 0)

	for _, user := range users {

		member := make(map[string]interface{})
		member["id"] = user.UID
		member["name"] = user.Username
		member["ownership"] = user.OwnType
		fmt.Printf("member: %v", member)

		members = append(members, member)
	}

	finalUser["users"] = members
	fmt.Printf("typeOf: %v", reflect.TypeOf(members))
	fmt.Printf("finalUser: %v", finalUser)
	finalUsersJSON, _ := json.Marshal(finalUser)

	return finalUsersJSON
}

func (j *Join) DeleteInvitation(tx *sqlx.Tx, inviteID int64) error {

	_, err := j.DeleteById(tx, inviteID)

	return err
}

func buildAddInviteResponseJSON(res sql.Result) []byte {

	inviteID, _ := res.LastInsertId()

	result := make(map[string]interface{})

	result["invite_id"] = inviteID
	result["message"] = "waiting for user to accept invitation"

	finalInviteJSON, _ := json.Marshal(result)
	return finalInviteJSON
}

func buildAddJoinResponseJSON(res sql.Result) []byte {

	inviteID, _ := res.LastInsertId()

	result := make(map[string]interface{})

	result["invite_id"] = inviteID
	result["message"] = "waiting for owner to accept join request"

	finalInviteJSON, _ := json.Marshal(result)
	return finalInviteJSON
}
