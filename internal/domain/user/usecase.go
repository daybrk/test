package user

import (
	"log/slog"
	"test-task/pkg/errs"
)

type Service interface {
	Validation(user *User, enrichmentUser *EnrichmentUser) bool
	EnrichmentAPI(user User) (EnrichmentUser, error)
	PutInDatabase(enrichmentUser EnrichmentUser) error
	CheckUserExist(id int) (bool, error)
	DeleteUser(id int) error
	ModifyUser(enrichmentUser EnrichmentUser) error
	GetFilteredUsers(filter Filter) ([]EnrichmentUser, error)
}

type userUseCase struct {
	service Service
	log     *slog.Logger
}

func NewUserUseCase(service Service, log *slog.Logger) *userUseCase {
	return &userUseCase{service: service, log: log}
}

func (uc userUseCase) Enrichment(user User) error {
	uc.log.Info("Происходит валидация")
	if uc.service.Validation(&user, nil) == false {
		uc.log.Warn("Некоректные входные данные")

		return errs.FioFailedErr
	}
	uc.log.Info("Валидация прошла успешно")

	uc.log.Info("Попытка обогатить данные пользователя")
	enrichmentUser, err := uc.service.EnrichmentAPI(user)
	if err != nil {
		uc.log.Error("ошибка при обогащении пользователя", slog.String("errs", err.Error()))

		return err
	}
	uc.log.Info("Обогащение успешно прошло", slog.Any("user", enrichmentUser))

	uc.log.Info("Попытка добавить пользователя в бд")
	err = uc.service.PutInDatabase(enrichmentUser)
	if err != nil {
		uc.log.Error("ошибка при добавлении пользователя", slog.String("errs", err.Error()))

		return err
	}
	uc.log.Info("Добавление в базу прошло успешно")

	uc.log.Info("Enrichment отработал без ошибок")

	return nil
}

func (uc userUseCase) DeleteUser(id int) error {
	uc.log.Info("Проверка пользователя на существование")
	if exist, err := uc.service.CheckUserExist(id); err == nil {
		if !exist {
			uc.log.Warn("пользователь не существует")

			return errs.UserNotFoundErr
		}
	} else {
		uc.log.Error("ошибка при проверке пользователя на существование", slog.String("errs", err.Error()))

		return err
	}
	uc.log.Info("Пользователь существует")

	uc.log.Info("Попытка удалить пользователя")
	err := uc.service.DeleteUser(id)
	if err != nil {
		uc.log.Error("не удалось удалить пользователя",
			slog.String("errs", err.Error()), slog.Int("id", id))

		return err
	}
	uc.log.Info("пользователь успешно удалён", slog.Int("id", id))

	uc.log.Info("DeleteUser отработал без ошибок")

	return nil
}

func (uc userUseCase) ModifyUser(enrichmentUser EnrichmentUser) error {
	uc.log.Info("Проверка пользователя на существование")
	if exist, err := uc.service.CheckUserExist(enrichmentUser.Id); err == nil {
		if !exist {
			uc.log.Warn("пользователь не существует")

			return errs.UserNotFoundErr
		}
	} else {
		uc.log.Error("ошибка при проверке пользователя на существование", slog.String("errs", err.Error()))

		return err
	}
	uc.log.Info("Проверка прошла успешно")

	uc.log.Info("Происходит валидация")
	if uc.service.Validation(nil, &enrichmentUser) == false {
		uc.log.Warn("Неправильные данные")

		return errs.FioFailedErr
	}
	uc.log.Info("Валидация прошла успешно")

	uc.log.Info("Попытка изменить пользователя")
	err := uc.service.ModifyUser(enrichmentUser)
	if err != nil {
		uc.log.Error("не удалось изменить пользователя",
			slog.String("errs", err.Error()), slog.Int("id", enrichmentUser.Id))
		return err
	}
	uc.log.Info("пользователь успешно изменён", slog.Int("id", enrichmentUser.Id))

	uc.log.Info("ModifyUser отработал без ошибок")

	return nil
}

func (uc userUseCase) Filter(filter Filter) ([]EnrichmentUser, error) {
	uc.log.Info("Попытка взять отфильтрованных пользователей")
	users, err := uc.service.GetFilteredUsers(filter)
	if err != nil {
		uc.log.Error("ошибка при попытке взять пользователей с использованием фильтров",
			slog.String("err", err.Error()))

		return nil, err
	}
	uc.log.Info("Пользователи успешно отфильтрованы и взяты")

	uc.log.Info("Filter отработал без ошибок")

	return users, nil
}
