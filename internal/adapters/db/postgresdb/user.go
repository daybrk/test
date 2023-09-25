package postgresdb

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log/slog"
	"strings"
	"test-task/pkg/errs"
)

type storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewUserStorage(db *sqlx.DB, log *slog.Logger) *storage {
	return &storage{db: db, log: log}
}

func (s storage) Insert(user EnrichmentUser) error {
	nationalityArray := "{" + strings.Join(user.Nationality, ",") + "}"

	_, err := s.db.Exec(`
		INSERT INTO public.user(name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, nationalityArray)
	if err != nil {
		s.log.Error("ошибка при выполнении запроса", slog.String("err", err.Error()))

		return err
	}

	return nil
}

func (s storage) DeleteUser(id int) error {
	ex, err := s.db.Exec(`
		DELETE FROM public.user WHERE id = $1`, id)
	if err != nil {
		s.log.Error("ошибка при выполнении запроса", slog.String("err", err.Error()))

		return err
	}

	aff, err := ex.RowsAffected()

	if aff == 0 {
		s.log.Warn("Пользователь с таким id не существует")

		return errs.UserNotFoundErr
	}

	return nil
}

func (s storage) UserExist(id int) error {
	r, err := s.db.Exec(`
		Select 1 FROM public.user WHERE id = $1`, id)
	if err != nil {
		s.log.Error("ошибка при выполнении запроса", slog.String("err", err.Error()))

		return err
	}

	aff, err := r.RowsAffected()

	if aff == 0 {
		s.log.Warn("Пользователь с таким id не существует")

		return sql.ErrNoRows
	}

	return nil
}

func (s storage) EditUser(user EnrichmentUser) error {
	nationalityArray := "{" + strings.Join(user.Nationality, ",") + "}"

	_, err := s.db.Exec(`
		UPDATE  public.user Set name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
		WHERE id = $7`, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, nationalityArray, user.Id)
	if err != nil {
		s.log.Error("ошибка при выполнении запроса", slog.String("err", err.Error()))

		return err
	}

	return nil
}

func (s storage) FilteredUsers(filter Filter) ([]EnrichmentUser, error) {
	var nationalityArray *string
	if filter.Nationality != nil {
		nationalityStr := "{" + strings.Join(filter.Nationality, ",") + "}"
		nationalityArray = &nationalityStr
	}

	r, err := s.db.Query(`
		Select * 
		FROM public.user 
		WHERE (name = $1 OR $1 IS NULL)
		  AND (surname = $2 OR $2 IS NULL)
		  AND (patronymic = $3 OR $3 IS NULL)
		  AND (age = $4 OR $4 IS NULL)
		  AND (gender = $5 OR $5 IS NULL)
		  AND (nationality = $6 OR $6 IS NULL)`,
		filter.Name, filter.Surname,
		filter.Patronymic, filter.Age,
		filter.Gender, nationalityArray)
	if err != nil {
		s.log.Error("ошибка при выполнении запроса", slog.String("err", err.Error()))

		return nil, err
	}

	var users []EnrichmentUser
	for r.Next() {
		var user EnrichmentUser
		err = r.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic,
			&user.Age, &user.Gender, (*pq.StringArray)(&user.Nationality))
		if err != nil {
			s.log.Error("ошибка при сканировании пользователя", slog.String("err", err.Error()))

			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
