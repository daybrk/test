package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"test-task/internal/adapters/web"
)

type Storage interface {
	Insert() error
	DeleteUser(id int) error
	UserExist(id int) error
	EditUser() error
}

type Router interface {
	EnrichmentAge(name string) (web.Age, error)
	EnrichmentGender(name string) (web.Gender, error)
	EnrichmentNationality(name string) (web.Nationality, error)
}

type userService struct {
	storage Storage
	router  Router
	log     *slog.Logger
}

func NewUserService(storage Storage, router Router, log *slog.Logger) *userService {
	fmt.Println("NewUserService")
	return &userService{storage: storage, router: router, log: log}
}

// TODO: добавить логов
func (e userService) Validation(fio *User, enrichmentFio *EnrichmentUser) bool {
	if fio != nil {
		if fio.Name == "" || hasCyrillic(fio.Name) {
			return false
		}

		if fio.Surname == "" || hasCyrillic(fio.Surname) {
			return false
		}
	}

	if enrichmentFio != nil {
		if enrichmentFio.Name == "" || hasCyrillic(enrichmentFio.Name) {
			return false
		}

		if enrichmentFio.Surname == "" || hasCyrillic(enrichmentFio.Surname) {
			return false
		}

		if enrichmentFio.Gender == "" || hasCyrillic(enrichmentFio.Gender) {
			return false
		}

		if enrichmentFio.Nationality == nil {
			return false
		}
	}

	return true
}

// TODO: Дополнить логи
func (e userService) EnrichmentAPI(fio User) (EnrichmentUser, error) {
	var enrichmentFIO EnrichmentUser
	enrichmentFIO.Name = fio.Name
	enrichmentFIO.Surname = fio.Surname
	enrichmentFIO.Patronymic = fio.Patronymic

	age, err := e.router.EnrichmentAge(fio.Name)
	if err != nil {
		return EnrichmentUser{}, err
	}
	enrichmentFIO.Age = age.Age
	e.log.Info("обогащение возрастом",
		slog.Any("age", age.Age), slog.Any("после обогащения", enrichmentFIO))

	gender, err := e.router.EnrichmentGender(fio.Name)
	if err != nil {
		return EnrichmentUser{}, err
	}
	enrichmentFIO.Gender = gender.Gender
	e.log.Info("обогащение полом",
		slog.Any("gender", gender.Gender), slog.Any("после обогащения", enrichmentFIO))

	nationality, err := e.router.EnrichmentNationality(fio.Name)
	if err != nil {
		return EnrichmentUser{}, err
	}
	for _, value := range nationality.Country {
		enrichmentFIO.Nationality = append(enrichmentFIO.Nationality, value.CountryId)
	}
	e.log.Info("обогащение национальностью",
		slog.Any("nationality", nationality.Country), slog.Any("после обогащения", enrichmentFIO))

	return enrichmentFIO, nil
}

func (e userService) CheckUserExist(id int) (bool, error) {
	err := e.storage.UserExist(id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		e.log.Warn("ошибка при получении пользователя", slog.String("err", err.Error()))

		return false, err
	}

	return true, nil
}

func (e userService) DeleteUser(id int) error {
	err := e.storage.DeleteUser(id)
	if err != nil {
		e.log.Error("ошибка при удалении пользователя",
			slog.String("err", err.Error()), slog.Int("user_id", id))

		return err
	}

	e.log.Info("пользователь успешно удалён", slog.Int("user_id", id))

	return nil
}

func (e userService) ModifyUser(enrichmentFio EnrichmentUser) error {
	err := e.storage.EditUser()
	if err != nil {
		e.log.Error("ошибка при изменении пользователя",
			slog.String("err", err.Error()), slog.Int("user_id", enrichmentFio.Id))

		return err
	}

	e.log.Info("пользователь успешно изменён", slog.Int("user_id", enrichmentFio.Id))

	return nil
}

func (e userService) PutInDatabase(enrichmentFio EnrichmentUser) error {
	err := e.storage.Insert()
	if err != nil {
		e.log.Error("ошибка при добавлении пользователя",
			slog.String("err", err.Error()), slog.Any("user", enrichmentFio))

		return err
	}

	e.log.Info("пользователь успешно добавлен", slog.Any("user", enrichmentFio))

	return nil
}

func hasCyrillic(input string) bool {
	re := regexp.MustCompile("[\u0400-\u04FF]+")

	return re.MatchString(input)
}
