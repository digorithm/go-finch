package models

type UserOwnType struct {
	Id			int64 `db: "id"`
	Email 		string `db: "email"`
	Password	string `db: "password"`
	Username	string `db: "username"`
	OwnType		int64  `db: "own_type"`
	Description string `db: "description"`
}

