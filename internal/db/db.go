package db

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file:data.db?cache=shared&mode=rwc")
	if err != nil {
		panic(err)
	}
	return &Database{
		db: sqldb,
	}
}

func (d Database) OpenDatabase() *bun.DB {

	db := bun.NewDB(d.db, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}
