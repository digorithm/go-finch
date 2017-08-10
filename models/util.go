package models

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestSetup(t *testing.T, db *sqlx.DB) func(t *testing.T, db *sqlx.DB) {
	db.SetMaxIdleConns(0)
	return func(t *testing.T, db *sqlx.DB) {
		if db.Stats().OpenConnections != 0 {
			t.Fatal("DB connections not zero")
		}
		db.Close()
	}
}
