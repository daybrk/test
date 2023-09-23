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

func NewUserStorage(db *sqlx.DB, log *slog.Logger) *storage {
	fmt.Println("NewUserStorage")
	return &storage{db: db, log: log}
}

func (s storage) Insert() error {
	//TODO implement me
	return nil
}

func (s storage) DeleteUser(id int) error {
	//TODO implement me
	return nil
}

func (s storage) UserExist(id int) error {
	//TODO implement me
	return nil
}

func (s storage) EditUser() error {
	//TODO implement me
	return nil
}
