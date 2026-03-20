package database

import (
	"github.com/hunderaweke/sma-go/domain"

	"gorm.io/gorm"
)

const (
	SQLite = iota
	Postgres
)

func NewDB(dbType int) (*gorm.DB, error) {
	switch dbType {
	case SQLite:
		return NewSQLiteDB("test.db")
	case Postgres:
		return NewPostgresConn()
	default:
		return nil, domain.InternalError(nil, "invalid db type")
	}
}
