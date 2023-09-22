package postgresdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewEnrichmentStorage(db *sqlx.DB, log *slog.Logger) *storage {
	fmt.Println("NewEnrichmentStorage")
	return &storage{db: db, log: log}
}

func (s storage) Insert() error {
	//TODO implement me
	panic("implement me")
}
