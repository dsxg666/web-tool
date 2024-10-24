package db

import "database/sql"

type Database struct {
	DbHandle *sql.DB
}
