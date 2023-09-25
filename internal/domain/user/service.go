package user

import (
	"database/sql"
	"errors"
	"log/slog"
	"regexp"
	"test-task/internal/adapters/db/postgresdb"
	"test-task/internal/adapters/web"
)

type Storage interface {
	Insert(user postgresdb.EnrichmentUser) error
	DeleteUser(id int) error
	UserExist(id int) error
	EditUser(user postgresdb.EnrichmentUser) error
	FilteredUsers(filter postgresdb.Filter) ([]postgresdb.EnrichmentUser, error)
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
	return &userService{storage: storage, router: router, log: log}
}

func (e userService) Validation(fio *User, enrichmentFio *EnrichmentUser) bool {
	if fio != nil {
		if fio.Name == "" || hasCyrillic(fio.Name) {
			e.log.Warn("Неправильно заполненное поле Name", slog.String("name", fio.Name))

			return false
		}

		if fio.Surname == "" || hasCyrillic(fio.Surname) {
			e.log.Warn("Неправильно заполненное поле Surname", slog.String("name", fio.Name))

			return false
		}
	}

	if enrichmentFio != nil {
		if enrichmentFio.Name == "" || hasCyrillic(enrichmentFio.Name) {
			e.log.Warn("Неправильно заполненное поле Name", slog.String("name", enrichmentFio.Name))

			return false
		}

		if enrichmentFio.Surname == "" || hasCyrillic(enrichmentFio.Surname) {
			e.log.Warn("Неправильно заполненное поле Surname", slog.String("surname", enrichmentFio.Name))

			return false
		}

		if enrichmentFio.Gender == "" || hasCyrillic(enrichmentFio.Gender) {
			e.log.Warn("Неправильно заполненное поле Gender", slog.String("gender", enrichmentFio.Gender))

			return false
		}

		if enrichmentFio.Nationality == nil {
			e.log.Warn("Неправильно заполненное поле Nationality",
				slog.Any("nationality", enrichmentFio.Nationality))

			return false
		}
	}

	return true
}

func (e userService) EnrichmentAPI(fio User) (EnrichmentUser, error) {
	var enrichmentFIO EnrichmentUser
	enrichmentFIO.Name = fio.Name
	enrichmentFIO.Surname = fio.Surname
	enrichmentFIO.Patronymic = fio.Patronymic

	e.log.Info("Попытка обратиться к api для получения age", slog.String("name", fio.Name))
	age, err := e.router.EnrichmentAge(fio.Name)
	if err != nil {
		e.log.Error("при использовании api для обогащения age произошла ошибка",
			slog.String("err", err.Error()))

		return EnrichmentUser{}, err
	}
	enrichmentFIO.Age = age.Age
	e.log.Info("обогащение возрастом",
		slog.Any("age", age.Age), slog.Any("после обогащения", enrichmentFIO))

	e.log.Info("Попытка обратиться к api для получения gender", slog.String("name", fio.Name))
	gender, err := e.router.EnrichmentGender(fio.Name)
	if err != nil {
		e.log.Error("при использовании api для обогащения gender произошла ошибка",
			slog.String("err", err.Error()))

		return EnrichmentUser{}, err
	}
	enrichmentFIO.Gender = gender.Gender
	e.log.Info("обогащение полом",
		slog.Any("gender", gender.Gender), slog.Any("после обогащения", enrichmentFIO))

	e.log.Info("Попытка обратиться к api для получения nationality", slog.String("name", fio.Name))
	nationality, err := e.router.EnrichmentNationality(fio.Name)
	if err != nil {
		e.log.Error("при использовании api для обогащения nationality произошла ошибка",
			slog.String("err", err.Error()))

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
	e.log.Info("Попытка обратиться к бд для проверки пользователя на существование", slog.Int("id", id))
	err := e.storage.UserExist(id)
	if errors.Is(err, sql.ErrNoRows) {
		e.log.Warn("пользователь не существует", slog.String("errs", err.Error()))

		return false, nil
	}
	if err != nil {
		e.log.Error("ошибка при получении пользователя", slog.String("errs", err.Error()))

		return false, err
	}
	e.log.Info("Пользователь существует", slog.Int("id", id))

	return true, nil
}

func (e userService) DeleteUser(id int) error {
	e.log.Info("Попытка обратиться к бд для удаление пользователя", slog.Int("id", id))
	err := e.storage.DeleteUser(id)
	if err != nil {
		e.log.Error("ошибка при удалении пользователя",
			slog.String("errs", err.Error()), slog.Int("id", id))

		return err
	}
	e.log.Info("пользователь успешно удалён", slog.Int("id", id))

	return nil
}

func (e userService) ModifyUser(enrichmentFio EnrichmentUser) error {
	e.log.Info("Попытка обратиться к бд для изменения пользователя", slog.Any("user_edit", enrichmentFio))
	err := e.storage.EditUser(postgresdb.EnrichmentUser{
		Id:          enrichmentFio.Id,
		Name:        enrichmentFio.Name,
		Surname:     enrichmentFio.Surname,
		Patronymic:  enrichmentFio.Patronymic,
		Age:         enrichmentFio.Age,
		Gender:      enrichmentFio.Gender,
		Nationality: enrichmentFio.Nationality,
	})
	if err != nil {
		e.log.Error("ошибка при изменении пользователя",
			slog.String("errs", err.Error()), slog.Int("user_id", enrichmentFio.Id))

		return err
	}
	e.log.Info("пользователь успешно изменён", slog.Int("user_id", enrichmentFio.Id))

	return nil
}

func (e userService) GetFilteredUsers(filter Filter) ([]EnrichmentUser, error) {
	e.log.Info("Попытка обратиться к бд для получение пользователей по фильтрам", slog.Any("filter", filter))
	users, err := e.storage.FilteredUsers(postgresdb.Filter{
		Name:        filter.Name,
		Surname:     filter.Surname,
		Patronymic:  filter.Patronymic,
		Age:         filter.Age,
		Gender:      filter.Gender,
		Nationality: filter.Nationality,
	})
	if err != nil {
		e.log.Error("ошибка при попытке взять пользователей с использованием фильтров",
			slog.String("err", err.Error()))

		return nil, err
	}
	e.log.Info("пользователи успешно получены", slog.Any("user", users))

	var us []EnrichmentUser
	for _, value := range users {
		us = append(us, EnrichmentUser{
			Id:          value.Id,
			Name:        value.Name,
			Surname:     value.Surname,
			Patronymic:  value.Patronymic,
			Age:         value.Age,
			Gender:      value.Gender,
			Nationality: value.Nationality,
		})
	}

	return us, nil
}

func (e userService) PutInDatabase(enrichmentFio EnrichmentUser) error {
	e.log.Info("Попытка обратиться к бд для добавления пользователя", slog.Any("user", enrichmentFio))
	err := e.storage.Insert(postgresdb.EnrichmentUser{
		Name:        enrichmentFio.Name,
		Surname:     enrichmentFio.Surname,
		Patronymic:  enrichmentFio.Patronymic,
		Age:         enrichmentFio.Age,
		Gender:      enrichmentFio.Gender,
		Nationality: enrichmentFio.Nationality,
	})
	if err != nil {
		e.log.Error("ошибка при добавлении пользователя",
			slog.String("errs", err.Error()), slog.Any("user", enrichmentFio))

		return err
	}
	e.log.Info("пользователь успешно добавлен", slog.Any("user", enrichmentFio))

	return nil
}

func hasCyrillic(input string) bool {
	re := regexp.MustCompile("[\u0400-\u04FF]+")

	return re.MatchString(input)
}
