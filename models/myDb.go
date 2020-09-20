package models

import "database/sql"

type MyDb struct {
	Db *sql.DB
	Tx *sql.Tx
}
