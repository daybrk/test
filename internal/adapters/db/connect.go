package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"test-task/internal/config"
)

func ConnectToPostgres(l *slog.Logger) (conn *sqlx.DB, err error) {
	conn, err = sqlx.Open("postgres", config.GetConfig().Postgres.ConnStr)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}
	l.Info("Подключение к базе успешно выполнено")

	_, err = conn.Exec(`
		Create Table IF NOT EXISTS public.user(
		    id SERIAL PRIMARY KEY,
		  	name VARCHAR NOT NULL,
			surname     VARCHAR NOT NULL,
			patronymic  VARCHAR,
			age         INTEGER,
			gender      VARCHAR,
			nationality TEXT[]
		)`)
	if err != nil {
		l.Error("Неполучилось провести миграцию", slog.String("err", err.Error()))

		return nil, err
	}
	l.Info("Миграция успешно выполнена")

	return conn, err
}
