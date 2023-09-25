package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"test-task/internal/config"
)

func ConnectToPostgres() (conn *sqlx.DB, err error) {
	conn, err = sqlx.Open("postgres", config.GetConfig().Postgres.ConnStr)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	_, err = conn.Exec(`
		Create Table IF NOT EXISTS main.user(
		    id SERIAL PRIMARY KEY,
		  	name VARCHAR NOT NULL,
			surname     VARCHAR NOT NULL,
			patronymic  VARCHAR,
			age         INTEGER,
			gender      VARCHAR,
			nationality TEXT[]
		)`)

	if err != nil {
		fmt.Println(err, "SADASDAS")

		return nil, err
	}

	return conn, err
}
