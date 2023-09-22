package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"test-task/internal/config"
)

var (
	Connection *sqlx.DB
)

func checkConnection() (err error) {
	if Connection.Ping() != nil {
		fmt.Println("Соединение прервано")
		err = ConnectToPostgres()
		if err != nil {
			fmt.Println(err)

			return err
		}
	}

	return nil
}

func ConnectToPostgres() (err error) {
	Connection, err = sqlx.Open("postgres", config.GetConfig().Postgres.ConnStr)
	if err != nil {
		fmt.Println(err)

		return err
	}

	return err
}
