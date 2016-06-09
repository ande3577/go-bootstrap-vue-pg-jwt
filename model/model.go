package model

import (
	"database/sql"
	"gopkg.in/gorp.v1"
)

var dbConnection *sql.DB
var dbMap *gorp.DbMap

func Initialize(db *sql.DB) {
	dbConnection = db
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
}
