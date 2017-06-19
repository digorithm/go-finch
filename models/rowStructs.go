package models


type HouseRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type UserRow struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
}

type OwnerRow struct {
	OwnType		int64  `db:"own_type"`
	Description string `db:"description"`
}

type UserOwnTypeRow struct {
	UserRow
	OwnerRow
}

type HouseStorageRow struct {
	HouseRow
	OwnerRow
}

